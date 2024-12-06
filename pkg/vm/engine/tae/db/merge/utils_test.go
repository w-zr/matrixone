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

package merge

import (
	"context"
	"github.com/matrixorigin/matrixone/pkg/defines"
	"github.com/matrixorigin/matrixone/pkg/fileservice"
	"github.com/matrixorigin/matrixone/pkg/pb/api"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"os"
	"path"
	"testing"
)

func TestResourceController(t *testing.T) {
	rc := new(resourceController)
	rc.setMemLimit(10000)
	require.Equal(t, int64(7500), rc.limit)
	require.Equal(t, int64(7500), rc.availableMem())

	rc.refresh()
	rc.limit = rc.using + 1
	require.Equal(t, int64(1), rc.availableMem())

	require.Panics(t, func() { rc.setMemLimit(0) })
}

func Test_CleanUpUselessFiles(t *testing.T) {
	tDir := os.TempDir()
	dir := path.Join(tDir, "/local")
	assert.NoError(t, os.RemoveAll(dir))
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	c := fileservice.Config{
		Name:    defines.ETLFileServiceName,
		Backend: "DISK",
		DataDir: dir,
		Cache:   fileservice.DisabledCacheConfig,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs, err := fileservice.NewFileService(ctx, c, nil)
	assert.Nil(t, err)
	defer fs.Close(ctx)

	ent := &api.MergeCommitEntry{
		BookingLoc: []string{"abc"},
	}

	CleanUpUselessFiles(ent, fs)
}

func TestScore(t *testing.T) {
	o1 := newSortedTestObjectEntry(t, 0, 3, 1)
	o2 := newSortedTestObjectEntry(t, 1, 4, 2)
	o3 := newSortedTestObjectEntry(t, 1, 5, 4)
	o4 := newSortedTestObjectEntry(t, 1, 100, 4)
	o5 := newSortedTestObjectEntry(t, 5, 10, math.MaxInt32)

	// should merge
	require.Less(t, 1.1, score([]*catalog.ObjectEntry{o1, o1}))
	require.Less(t, 1.1, score([]*catalog.ObjectEntry{o1, o2}))
	require.Less(t, 1.1, score([]*catalog.ObjectEntry{o1, o3}))
	// should not merge
	require.Greater(t, 1.1, score([]*catalog.ObjectEntry{o1, o4}))
	require.Greater(t, 1.1, score([]*catalog.ObjectEntry{o1, o2, o4}))
	require.Greater(t, 1.1, score([]*catalog.ObjectEntry{o1, o5}))

	o6 := newTestVarcharObjectEntry(t, "a", "z", 1)
	o7 := newTestVarcharObjectEntry(t, "b", "y", 1)

	require.Less(t, 1.1, score([]*catalog.ObjectEntry{o6, o7}))
}

func TestRemoveOversize(t *testing.T) {
	o1 := newSortedTestObjectEntry(t, 0, 0, 1)
	o2 := newSortedTestObjectEntry(t, 0, 0, 2)
	o3 := newSortedTestObjectEntry(t, 0, 0, 4)
	o5 := newSortedTestObjectEntry(t, 0, 0, math.MaxInt32)

	require.ElementsMatch(t, []*catalog.ObjectEntry{o1, o2}, removeOversize([]*catalog.ObjectEntry{o1, o2}))
	require.ElementsMatch(t, []*catalog.ObjectEntry{o1, o2}, removeOversize([]*catalog.ObjectEntry{o5, o1, o2}))
	require.ElementsMatch(t, nil, removeOversize([]*catalog.ObjectEntry{o1, o3}))
}

func BenchmarkRemoveOversize(b *testing.B) {
	o1 := newSortedTestObjectEntry(b, 0, 50, math.MaxInt32)
	o2 := newSortedTestObjectEntry(b, 51, 100, 1)
	o3 := newSortedTestObjectEntry(b, 49, 52, 2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		removeOversize([]*catalog.ObjectEntry{o1, o2, o3})
	}
}
