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

package malloc

import (
	"runtime"
	"testing"
)

func TestMakeAllocator(t *testing.T) {
	testAllocator(t, func() Allocator {
		return NewShardedAllocator(
			runtime.GOMAXPROCS(0),
			func() Allocator {
				return NewClassAllocator(NewFixedSizeMakeAllocator)
			},
		)
	})
}

func BenchmarkMakeAllocator(b *testing.B) {
	for _, n := range benchNs {
		benchmarkAllocator(b, func() Allocator {
			return NewShardedAllocator(
				runtime.GOMAXPROCS(0),
				func() Allocator {
					return NewClassAllocator(NewFixedSizeMakeAllocator)
				},
			)
		}, n)
	}
}

func FuzzMakeAllocator(f *testing.F) {
	fuzzAllocator(f, func() Allocator {
		return NewShardedAllocator(
			runtime.GOMAXPROCS(0),
			func() Allocator {
				return NewClassAllocator(NewFixedSizeMakeAllocator)
			},
		)
	})
}
