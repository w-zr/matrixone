package tasks

import (
	"hash/fnv"
	"sync/atomic"

	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/common"
)

var DefaultScopeHasher = func(scope *common.ID) int {
	if scope == nil {
		return 0
	}
	hasher := fnv.New64a()
	hasher.Write(types.EncodeUint64(&scope.TableID))
	hasher.Write(types.EncodeUuid(scope.SegmentID()))
	return int(hasher.Sum64())
}

type scopeHasher func(scope *common.ID) int
type ShardedTaskHandler struct {
	handlers []TaskHandler
	hasher   scopeHasher
	curr     atomic.Uint64
}

func NewShardedTaskHandler(hasher scopeHasher) *ShardedTaskHandler {
	d := &ShardedTaskHandler{
		handlers: make([]TaskHandler, 0),
	}
	if hasher == nil {
		d.hasher = d.roundRobinSharder
	} else {
		d.hasher = hasher
	}
	return d
}

func (d *ShardedTaskHandler) AddHandle(h TaskHandler) {
	d.handlers = append(d.handlers, h)
}

func (d *ShardedTaskHandler) roundRobinSharder(scope *common.ID) int {
	return int(d.curr.Add(1))
}

func (d *ShardedTaskHandler) Enqueue(task Task) {
	shardIdx := d.hasher(task.(ScopedTask).Scope()) % len(d.handlers)
	d.handlers[shardIdx].Enqueue(task)
}

func (d *ShardedTaskHandler) Start() {
	for _, h := range d.handlers {
		h.Start()
	}
}

func (d *ShardedTaskHandler) Close() error {
	for _, h := range d.handlers {
		h.Close()
	}
	return nil
}
