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

package tasks

import (
	"context"
	"hash/fnv"
	"sync/atomic"
	"time"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/common"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/iface/txnif"
)

var (
	ErrBadTaskRequestPara    = moerr.NewInternalErrorNoCtx("tae scheduler: bad task request parameters")
	ErrScheduleScopeConflict = moerr.NewInternalErrorNoCtx("tae scheduler: scope conflict")
)

type FuncT func() error

type TaskType uint16

var taskIdAllocator *common.IdAllocator

const (
	NoopTask TaskType = iota
	MockTask
	CustomizedTask

	DataCompactionTask
	CheckpointTask
	GCTask
	IOTask
	FlushTableTailTask
)

var taskNames = map[TaskType]string{
	NoopTask:           "Noop",
	MockTask:           "Mock",
	DataCompactionTask: "Compaction",
	CheckpointTask:     "Checkpoint",
	GCTask:             "GC",
	IOTask:             "IO",
	FlushTableTailTask: "FlushTableTail",
}

func RegisterType(t TaskType, name string) {
	_, ok := taskNames[t]
	if ok {
		panic(moerr.NewInternalErrorNoCtxf("duplicate task type: %d, %s", t, name))
	}
	taskNames[t] = name
}

func TaskName(t TaskType) string {
	return taskNames[t]
}

func init() {
	taskIdAllocator = common.NewIdAllocator(1)
}

type TxnTaskFactory = func(ctx *Config, txn txnif.AsyncTxn) (Task, error)

func NextTaskId() uint64 {
	return taskIdAllocator.Alloc()
}

type Task interface {
	OnExec(ctx context.Context) error
	SetError(err error)
	GetError() error
	WaitDone(ctx context.Context) error
	Waitable() bool
	GetCreateTime() time.Time
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetExecuteTime() int64
	AddObserver(Observer)
	ID() uint64
	Type() TaskType
	Cancel() error
	Name() string
}

type ScopedTask interface {
	Task
	Scope() *common.ID
}

type MScopedTask interface {
	Task
	Scopes() []common.ID
}

var DefaultScopeSharder = func(scope *common.ID) int {
	if scope == nil {
		return 0
	}
	hasher := fnv.New64a()
	hasher.Write(types.EncodeUint64(&scope.TableID))
	hasher.Write(types.EncodeUuid(scope.SegmentID()))
	return int(hasher.Sum64())
}

type FnTask struct {
	*BaseTask
	fn FuncT
}

func NewFnTask(ctx *Config, taskType TaskType, fn FuncT) *FnTask {
	task := &FnTask{
		fn: fn,
	}
	task.BaseTask = NewBaseTask(task, taskType, ctx)
	return task
}

func (task *FnTask) Execute(ctx context.Context) error {
	return task.fn()
}

type ScopedFnTask struct {
	*FnTask
	scope *common.ID
}

func NewScopedFnTask(ctx *Config, taskType TaskType, scope *common.ID, fn FuncT) *ScopedFnTask {
	task := &ScopedFnTask{
		FnTask: new(FnTask),
		scope:  scope,
	}
	task.fn = fn
	task.BaseTask = NewBaseTask(task, taskType, ctx)
	return task
}

func (task *ScopedFnTask) Scope() *common.ID { return task.scope }

type MultiScopedFnTask struct {
	*FnTask
	scopes []common.ID
}

func NewMultiScopedFnTask(ctx *Config, taskType TaskType, scopes []common.ID, fn FuncT) *MultiScopedFnTask {
	task := &MultiScopedFnTask{
		FnTask: new(FnTask),
		scopes: scopes,
	}
	task.fn = fn
	task.BaseTask = NewBaseTask(task, taskType, ctx)
	return task
}

func (task *MultiScopedFnTask) Scopes() []common.ID {
	return task.scopes
}

type Op struct {
	impl       IOpInternal
	errorC     chan error
	waitedOnce atomic.Bool
	err        error
	result     any
	createTime time.Time
	startTime  time.Time
	endTime    time.Time
	doneCB     opExecFunc
	observers  []Observer
}

type Observer interface {
	OnExecDone(any)
}

type IOpInternal interface {
	PreExecute() error
	Execute(ctx context.Context) error
	PostExecute() error
}

type taskHandle interface {
	OnExec()
	OnStopped()
}
