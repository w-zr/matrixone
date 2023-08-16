// Copyright 2021 Matrix Origin
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

package objectio

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/matrixorigin/matrixone/pkg/fileservice"
	"github.com/matrixorigin/matrixone/pkg/fileservice/lrucache"
	"github.com/matrixorigin/matrixone/pkg/util/toml"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/blockio"
)

type CacheConfig struct {
	MemoryCapacity toml.ByteSize `toml:"memory-capacity"`
}

var metaCache *lrucache.LRU[ObjectNameShort, fileservice.Bytes]
var bloomFilterCache *lrucache.LRU[ObjectNameShort, BloomFilter]
var onceInit sync.Once
var metaCacheStats hitStats
var metaCacheHitStats hitStats
var bloomFilterCacheStats hitStats

func init() {
	metaCache = lrucache.New[ObjectNameShort, fileservice.Bytes](512*1024*1024, nil, nil, nil)
	bloomFilterCache = lrucache.New[ObjectNameShort, BloomFilter](512*1024*1024, nil, nil, nil)
}

func InitMetaCache(size int64) {
	onceInit.Do(func() {
		metaCache = lrucache.New[ObjectNameShort, fileservice.Bytes](size, nil, nil, nil)
		bloomFilterCache = lrucache.New[ObjectNameShort, BloomFilter](size, nil, nil, nil)
	})
}

func ExportMetaCacheStats() string {
	var buf bytes.Buffer
	hw, hwt := metaCacheHitStats.ExportW()
	ht, htt := metaCacheHitStats.Export()
	w, wt := metaCacheStats.ExportW()
	t, tt := metaCacheStats.Export()
	fmt.Fprintf(
		&buf,
		"MetaCacheWindow: %d/%d | %d/%d, MetaCacheTotal: %d/%d | %d/%d",
		hw, hwt, w, wt, ht, htt, t, tt,
	)
	return buf.String()
}

func ExportBloomFilterCacheStats() string {
	var buf bytes.Buffer
	hw, hwt := bloomFilterCacheStats.ExportW()
	ht, htt := bloomFilterCacheStats.Export()
	fmt.Fprintf(
		&buf,
		"BloomFilterCacheWindow: %d/%d, BloomFilterCacheTotal: %d/%d",
		hw, hwt, ht, htt,
	)
	return buf.String()
}

func LoadObjectMetaByExtent(
	ctx context.Context,
	name *ObjectName,
	extent *Extent,
	prefetch bool,
	noLRUCache bool,
	fs fileservice.FileService,
) (meta ObjectMeta, err error) {
	v, ok := metaCache.Get(ctx, *name.Short(), false)
	if ok {
		var obj any
		obj, err = Decode(v)
		if err != nil {
			return
		}
		meta = obj.(ObjectMeta)
		metaCacheStats.Record(1, 1)
		if !prefetch {
			metaCacheHitStats.Record(1, 1)
		}
		return
	}
	if v, err = ReadExtent(ctx, name.String(), extent, noLRUCache, fs, constructorFactory); err != nil {
		return
	}
	var obj any
	obj, err = Decode(v)
	if err != nil {
		return
	}
	meta = obj.(ObjectMeta)
	metaCache.Set(ctx, *name.Short(), v[:], false)
	metaCacheStats.Record(0, 1)
	if !prefetch {
		metaCacheHitStats.Record(0, 1)
	}
	return
}

func FastLoadBF(
	ctx context.Context,
	loc Location,
	fs fileservice.FileService,
) (BloomFilter, error) {
	v, ok := bloomFilterCache.Get(ctx, *loc.Name().Short(), false)
	if ok {
		bloomFilterCacheStats.Record(1, 1)
		return v.Bytes(), nil
	}
	r, _ := blockio.NewObjectReader(fs, loc)
	v, _, err := r.LoadAllBF(ctx)
	if err != nil {
		return nil, err
	}
	bloomFilterCache.Set(ctx, *loc.ShortName(), v, false)
	bloomFilterCacheStats.Record(0, 1)
	return v, nil
}

func FastLoadObjectMeta(
	ctx context.Context,
	location *Location,
	prefetch bool,
	fs fileservice.FileService,
) (ObjectMeta, error) {
	extent := location.Extent()
	name := location.Name()
	return LoadObjectMetaByExtent(ctx, &name, &extent, prefetch, true, fs)
}
