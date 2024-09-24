// Copyright 2022 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/docker/go-units"
)

const (
	DefaultMinOsizeQualifiedMB   = 110   // MB
	DefaultMaxOsizeObjMB         = 128   // MB
	DefaultMinCNMergeSize        = 80000 // MB
	DefaultCNMergeMemControlHint = 8192  // MB
	DefaultMaxMergeObjN          = 16

	Const1GBytes = 1 << 30
	Const1MBytes = 1 << 20
)

var (
	RuntimeMaxMergeObjN        atomic.Int32
	RuntimeOsizeRowsQualified  atomic.Uint32
	RuntimeMaxObjOsize         atomic.Uint32
	RuntimeMinCNMergeSize      atomic.Uint64
	RuntimeCNMergeMemControl   atomic.Uint64
	RuntimeCNTakeOverAll       atomic.Bool
	IsStandaloneBoost          atomic.Bool
	ShouldStandaloneCNTakeOver atomic.Bool

	RuntimeOverallFlushMemCap atomic.Uint64

	FlushGapDuration atomic.Value
	FlushMemCapacity atomic.Int32
)

func init() {
	RuntimeMaxMergeObjN.Store(DefaultMaxMergeObjN)
	RuntimeOsizeRowsQualified.Store(DefaultMinOsizeQualifiedMB * Const1MBytes)
	RuntimeMaxObjOsize.Store(DefaultMaxOsizeObjMB * Const1MBytes)
	FlushGapDuration.Store(time.Minute)
	FlushMemCapacity.Store(20 * Const1MBytes)
}

///////
/// statistics component
///////

///
/// Table statistics
///

type TableCompactStat struct {
	sync.RWMutex

	// FlushDeadline is the deadline to flush table tail.
	FlushDeadline time.Time
	// LastMergeTime is the last merge time.
	LastMergeTime time.Time
}

func (s *TableCompactStat) ResetDeadline() {
	// add random +/- 10%
	s.Lock()
	defer s.Unlock()
	factor := 1.0 + float64(rand.Intn(21)-10)/100.0
	s.FlushDeadline = time.Now().Add(time.Duration(factor * float64(FlushGapDuration.Load().(time.Duration))))
}

func (s *TableCompactStat) AddMerge() {
	s.Lock()
	defer s.Unlock()
	s.LastMergeTime = time.Now()
}

////
// Other utils
////

func HumanReadableBytes(bytes int) string {
	return units.HumanSize(float64(bytes))
}
