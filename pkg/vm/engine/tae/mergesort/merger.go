// Copyright 2024 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mergesort

import (
	"context"
	"errors"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/container/batch"
	"github.com/matrixorigin/matrixone/pkg/container/nulls"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/pb/api"
	"github.com/matrixorigin/matrixone/pkg/sort"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/blockio"
)

type Merger interface {
	merge(context.Context) error
}

type releasableBatch struct {
	bat      *batch.Batch
	releaseF func()
}

type merger[T comparable] struct {
	heap *heapSlice[T]

	df      vector.DataSlice[T]
	deletes []*nulls.Nulls
	nulls   []*nulls.Nulls

	buffer *batch.Batch

	bats   []releasableBatch
	rowIdx []uint32

	objCnt           int
	objBlkCnts       []int
	accObjBlkCnts    []int
	loadedObjBlkCnts []int

	host MergeTaskHost

	writer *blockio.BlockWriter

	sortKeyIdx int

	rowPerBlk uint32
	stats     mergeStats
}

func newMerger[T comparable](host MergeTaskHost, lessFunc sort.LessFunc[T], sortKeyPos int, df vector.DataSlice[T]) Merger {
	size := host.GetObjectCnt()
	rowSizeU64 := host.GetTotalSize() / uint64(host.GetTotalRowCnt())
	m := &merger[T]{
		host:   host,
		objCnt: size,

		df:         df,
		bats:       make([]releasableBatch, size),
		rowIdx:     make([]uint32, size),
		deletes:    make([]*nulls.Nulls, size),
		nulls:      make([]*nulls.Nulls, size),
		heap:       newHeapSlice(size, lessFunc),
		sortKeyIdx: sortKeyPos,

		accObjBlkCnts: host.GetAccBlkCnts(),
		objBlkCnts:    host.GetBlkCnts(),
		rowPerBlk:     host.GetBlockMaxRows(),
		stats: mergeStats{
			totalRowCnt:   host.GetTotalRowCnt(),
			rowSize:       uint32(rowSizeU64),
			targetObjSize: host.GetTargetObjSize(),
			blkPerObj:     host.GetObjectMaxBlocks(),
		},
		loadedObjBlkCnts: make([]int, size),
	}
	totalBlkCnt := 0
	for _, cnt := range m.objBlkCnts {
		totalBlkCnt += cnt
	}
	if host.DoTransfer() {
		host.InitTransferMaps(totalBlkCnt)
	}

	return m
}

func (m *merger[T]) merge(ctx context.Context) error {
	for i := 0; i < m.objCnt; i++ {
		if ok, err := m.loadBlk(ctx, uint32(i)); !ok {
			if err == nil {
				continue
			}
			return errors.Join(moerr.NewInternalError(ctx, "failed to load first blk"), err)
		}

		heapPush(m.heap, heapElem[T]{
			data:   m.df.At(i, 0),
			isNull: m.nulls[i].Contains(0),
			src:    uint32(i),
		})
	}
	defer m.release()

	var releaseF func()
	for i := 0; i < m.objCnt; i++ {
		if m.bats[i].bat != nil {
			m.buffer, releaseF = getSimilarBatch(m.bats[i].bat, int(m.rowPerBlk), m.host)
			break
		}
	}
	// all batches are empty.
	if m.buffer == nil {
		return nil
	}
	defer releaseF()

	transferMaps := m.host.GetTransferMaps()
	for m.heap.Len() != 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		objIdx := m.nextPos()
		if m.deletes[objIdx].Contains(uint64(m.rowIdx[objIdx])) {
			// row is deleted
			if err := m.pushNewElem(ctx, objIdx); err != nil {
				return err
			}
			continue
		}
		rowIdx := m.rowIdx[objIdx]
		for i := range m.buffer.Vecs {
			err := m.buffer.Vecs[i].UnionOne(m.bats[objIdx].bat.Vecs[i], int64(rowIdx), m.host.GetMPool())
			if err != nil {
				return err
			}
		}

		if m.host.DoTransfer() {
			transferMaps[m.accObjBlkCnts[objIdx]+m.loadedObjBlkCnts[objIdx]-1][rowIdx] = api.TransferDestPos{
				ObjIdx: uint8(m.stats.objCnt),
				BlkIdx: uint16(m.stats.objBlkCnt),
				RowIdx: uint32(m.stats.blkRowCnt),
			}
		}

		m.stats.blkRowCnt++
		m.stats.objRowCnt++
		m.stats.mergedRowCnt++
		// write new block
		if m.stats.blkRowCnt == int(m.rowPerBlk) {
			m.stats.blkRowCnt = 0
			m.stats.objBlkCnt++

			if m.writer == nil {
				m.writer = m.host.PrepareNewWriter()
			}
			if _, err := m.writer.WriteBatch(m.buffer); err != nil {
				return err
			}
			// force clean
			m.buffer.CleanOnlyData()

			// write new object
			if m.stats.needNewObject() {
				// write object and reset writer
				if err := m.syncObject(ctx); err != nil {
					return err
				}
				// reset writer after sync
				m.stats.objBlkCnt = 0
				m.stats.objRowCnt = 0
				m.stats.objCnt++
			}
		}

		if err := m.pushNewElem(ctx, objIdx); err != nil {
			return err
		}
	}

	// write remain data
	if m.stats.blkRowCnt > 0 {
		m.stats.objBlkCnt++

		if m.writer == nil {
			m.writer = m.host.PrepareNewWriter()
		}
		if _, err := m.writer.WriteBatch(m.buffer); err != nil {
			return err
		}
		m.buffer.CleanOnlyData()
	}
	if m.stats.objBlkCnt > 0 {
		if err := m.syncObject(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *merger[T]) nextPos() uint32 {
	return heapPop[T](m.heap).src
}

func (m *merger[T]) loadBlk(ctx context.Context, objIdx uint32) (bool, error) {
	nextBatch, del, releaseF, err := m.host.LoadNextBatch(ctx, objIdx)
	if m.bats[objIdx].bat != nil {
		m.bats[objIdx].releaseF()
		m.bats[objIdx].releaseF = nil
	}
	if err != nil {
		if errors.Is(err, ErrNoMoreBlocks) {
			return false, nil
		}
		return false, err
	}
	for nextBatch.RowCount() == 0 {
		releaseF()
		nextBatch, del, releaseF, err = m.host.LoadNextBatch(ctx, objIdx)
		if err != nil {
			if errors.Is(err, ErrNoMoreBlocks) {
				return false, nil
			}
			return false, err
		}
	}

	m.bats[objIdx] = releasableBatch{bat: nextBatch, releaseF: releaseF}
	m.loadedObjBlkCnts[objIdx]++

	vec := nextBatch.GetVector(int32(m.sortKeyIdx))
	m.df.MustToCol(vec, objIdx)
	m.nulls[objIdx] = vec.GetNulls()
	m.deletes[objIdx] = del
	m.rowIdx[objIdx] = 0
	return true, nil
}

func (m *merger[T]) pushNewElem(ctx context.Context, objIdx uint32) error {
	m.rowIdx[objIdx]++
	if m.rowIdx[objIdx] >= uint32(m.df.Length(int(objIdx))) {
		if ok, err := m.loadBlk(ctx, objIdx); !ok {
			return err
		}
	}
	nextRow := m.rowIdx[objIdx]
	heapPush(m.heap, heapElem[T]{
		data:   m.df.At(int(objIdx), int(nextRow)),
		isNull: m.nulls[objIdx].Contains(uint64(nextRow)),
		src:    objIdx,
	})
	return nil
}

func (m *merger[T]) syncObject(ctx context.Context) error {
	if _, _, err := m.writer.Sync(ctx); err != nil {
		return err
	}
	cobjstats := m.writer.GetObjectStats()
	commitEntry := m.host.GetCommitEntry()
	commitEntry.CreatedObjs = append(commitEntry.CreatedObjs, cobjstats.Clone().Marshal())
	m.writer = nil
	return nil
}

func (m *merger[T]) release() {
	for _, bat := range m.bats {
		if bat.releaseF != nil {
			bat.releaseF()
		}
	}
}

func mergeObjs(ctx context.Context, mergeHost MergeTaskHost, sortKeyPos int) error {
	var merger Merger
	typ := mergeHost.GetSortKeyType()
	size := mergeHost.GetObjectCnt()
	if typ.IsVarlen() {
		df := vector.NewVarlenaDataSlice(size)
		merger = newMerger(mergeHost, sort.GenericLess[string], sortKeyPos, df)
	} else {
		switch typ.Oid {
		case types.T_bool:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[bool], size)
			merger = newMerger(mergeHost, sort.BoolLess, sortKeyPos, df)
		case types.T_bit:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[uint64], size)
			merger = newMerger(mergeHost, sort.GenericLess[uint64], sortKeyPos, df)
		case types.T_int8:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[int8], size)
			merger = newMerger(mergeHost, sort.GenericLess[int8], sortKeyPos, df)
		case types.T_int16:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[int16], size)
			merger = newMerger(mergeHost, sort.GenericLess[int16], sortKeyPos, df)
		case types.T_int32:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[int32], size)
			merger = newMerger(mergeHost, sort.GenericLess[int32], sortKeyPos, df)
		case types.T_int64:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[int64], size)
			merger = newMerger(mergeHost, sort.GenericLess[int64], sortKeyPos, df)
		case types.T_float32:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[float32], size)
			merger = newMerger(mergeHost, sort.GenericLess[float32], sortKeyPos, df)
		case types.T_float64:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[float64], size)
			merger = newMerger(mergeHost, sort.GenericLess[float64], sortKeyPos, df)
		case types.T_uint8:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[uint8], size)
			merger = newMerger(mergeHost, sort.GenericLess[uint8], sortKeyPos, df)
		case types.T_uint16:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[uint16], size)
			merger = newMerger(mergeHost, sort.GenericLess[uint16], sortKeyPos, df)
		case types.T_uint32:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[uint32], size)
			merger = newMerger(mergeHost, sort.GenericLess[uint32], sortKeyPos, df)
		case types.T_uint64:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[uint64], size)
			merger = newMerger(mergeHost, sort.GenericLess[uint64], sortKeyPos, df)
		case types.T_date:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Date], size)
			merger = newMerger(mergeHost, sort.GenericLess[types.Date], sortKeyPos, df)
		case types.T_timestamp:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Timestamp], size)
			merger = newMerger(mergeHost, sort.GenericLess[types.Timestamp], sortKeyPos, df)
		case types.T_datetime:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Datetime], size)
			merger = newMerger(mergeHost, sort.GenericLess[types.Datetime], sortKeyPos, df)
		case types.T_time:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Time], size)
			merger = newMerger(mergeHost, sort.GenericLess[types.Time], sortKeyPos, df)
		case types.T_enum:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Enum], size)
			merger = newMerger(mergeHost, sort.GenericLess[types.Enum], sortKeyPos, df)
		case types.T_decimal64:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Decimal64], size)
			merger = newMerger(mergeHost, sort.Decimal64Less, sortKeyPos, df)
		case types.T_decimal128:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Decimal128], size)
			merger = newMerger(mergeHost, sort.Decimal128Less, sortKeyPos, df)
		case types.T_uuid:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Uuid], size)
			merger = newMerger(mergeHost, sort.UuidLess, sortKeyPos, df)
		case types.T_TS:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.TS], size)
			merger = newMerger(mergeHost, sort.TsLess, sortKeyPos, df)
		case types.T_Rowid:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Rowid], size)
			merger = newMerger(mergeHost, sort.RowidLess, sortKeyPos, df)
		case types.T_Blockid:
			df := vector.NewFixedDataSlice(vector.MustFixedColNoTypeCheck[types.Blockid], size)
			merger = newMerger(mergeHost, sort.BlockidLess, sortKeyPos, df)
		default:
			return moerr.NewErrUnsupportedDataType(ctx, typ)
		}
	}
	return merger.merge(ctx)
}
