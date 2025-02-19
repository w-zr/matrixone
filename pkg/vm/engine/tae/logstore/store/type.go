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

package store

import (
	"context"

	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/logstore/driver"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/logstore/entry"
)

const (
	GroupCKP      = entry.GTCKp
	GroupInternal = entry.GTInternal
	GroupFiles    = entry.GTFiles
)

type Store interface {
	AppendEntry(gid uint32, entry entry.Entry) (lsn uint64, err error)

	// it always checkpoint the GroupPrepare group
	RangeCheckpoint(start, end uint64, files ...string) (ckpEntry entry.Entry, err error)

	// only used in the test
	// it returns the next lsn of group `GroupPrepare`
	GetCurrSeqNum() (lsn uint64)
	// only used in the test
	// it returns the remaining entries of group `GroupPrepare` to be checkpointed
	// GetCurrSeqNum() - GetCheckpointed()
	GetPendding() (cnt uint64)
	// only used in the test
	// it returns the last lsn of group `GroupPrepare` that has been checkpointed
	GetCheckpointed() (lsn uint64)
	// only used in the test
	// it returns the truncated dsn
	GetTruncated() uint64

	Replay(
		ctx context.Context,
		h ApplyHandle,
		modeGetter func() driver.ReplayMode,
		opt *driver.ReplayOption,
	) error

	Close() error
}

type ApplyHandle = func(group uint32, commitId uint64, payload []byte, typ uint16, info any) driver.ReplayEntryState
