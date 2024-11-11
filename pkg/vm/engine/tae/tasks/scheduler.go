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
	"fmt"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/options"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/logutil"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/common"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/txn/txnbase"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/wal"
)

var (
	ErrDispatcherNotFound = moerr.NewInternalErrorNoCtx("tae sched: dispatcher not found")
	ErrSchedule           = moerr.NewInternalErrorNoCtx("tae sched: cannot schedule")
)

type Scheduler interface {
	Start()
	Stop()
	Schedule(Task) error

	ScheduleTxnTask(ctx *Config, taskType TaskType, factory TxnTaskFactory) (Task, error)
	ScheduleMultiScopedTxnTask(ctx *Config, taskType TaskType, scopes []common.ID, factory TxnTaskFactory) (Task, error)
	ScheduleMultiScopedTxnTaskWithObserver(ctx *Config, taskType TaskType, scopes []common.ID, factory TxnTaskFactory, observers ...Observer) (Task, error)
	ScheduleMultiScopedFn(ctx *Config, taskType TaskType, scopes []common.ID, fn FuncT) (Task, error)
	ScheduleFn(ctx *Config, taskType TaskType, fn func() error) (Task, error)
	ScheduleScopedFn(ctx *Config, taskType TaskType, scope *common.ID, fn func() error) (Task, error)

	CheckAsyncScopes(scopes []common.ID) error

	GetCheckpointedLSN() uint64
	GetPenddingLSNCnt() uint64
}

type taskScheduler struct {
	typedHandler *TypedTaskHandler
	ctx          context.Context
	txnMgr       *txnbase.TxnManager
	wal          wal.Driver
}

func (s *taskScheduler) Start() {
	s.typedHandler.Start()
}

func NewTaskScheduler(ctx context.Context, txnMgr *txnbase.TxnManager, wal wal.Driver, cfg *options.SchedulerCfg) *taskScheduler {
	if cfg.AsyncWorkers < 0 {
		panic(fmt.Sprintf("bad param: %d txn workers", cfg.AsyncWorkers))
	}
	if cfg.IOWorkers < 0 {
		panic(fmt.Sprintf("bad param: %d io workers", cfg.IOWorkers))
	}
	s := &taskScheduler{
		ctx:          ctx,
		typedHandler: NewTypedTaskHandler(),
		txnMgr:       txnMgr,
		wal:          wal,
	}
	jobHandler := NewScopeCheckTaskHandlerDispatcher()
	mergeHandler := NewPoolHandler(ctx, cfg.AsyncWorkers)
	jobHandler.AddHandle(DataCompactionTask, mergeHandler)
	gcHandler := NewBaseTaskHandler(ctx, "gc")
	jobHandler.AddHandle(GCTask, gcHandler)
	flushHandler := NewPoolHandler(ctx, cfg.AsyncWorkers)
	jobHandler.AddHandle(FlushTableTailTask, flushHandler)

	ckpHandler := NewShardedTaskHandler(defaultScopeHasher)
	for i := 0; i < 4; i++ {
		handler := NewBaseTaskHandler(ctx, fmt.Sprintf("[ckpworker-%d]", i))
		ckpHandler.AddHandle(handler)
	}

	ioHandler := NewShardedTaskHandler(nil)
	for i := 0; i < cfg.IOWorkers; i++ {
		handler := NewBaseTaskHandler(ctx, fmt.Sprintf("[ioworker-%d]", i))
		ioHandler.AddHandle(handler)
	}

	s.typedHandler.AddHandle(GCTask, jobHandler)
	s.typedHandler.AddHandle(DataCompactionTask, jobHandler)
	s.typedHandler.AddHandle(IOTask, ioHandler)
	s.typedHandler.AddHandle(CheckpointTask, ckpHandler)
	s.Start()
	return s
}

func (s *taskScheduler) Stop() {
	s.typedHandler.Close()
	logutil.Info("TaskScheduler Stopped")
}

func (s *taskScheduler) ScheduleTxnTask(
	cfg *Config,
	taskType TaskType,
	factory TxnTaskFactory) (task Task, err error) {
	task = NewScheduledTxnTask(s.ctx, s.txnMgr, taskType, cfg, nil, factory)
	err = s.Schedule(task)
	return
}

func (s *taskScheduler) ScheduleMultiScopedTxnTask(
	cfg *Config,
	taskType TaskType,
	scopes []common.ID,
	factory TxnTaskFactory) (task Task, err error) {
	task = NewScheduledTxnTask(s.ctx, s.txnMgr, taskType, cfg, scopes, factory)
	err = s.Schedule(task)
	return
}

func (s *taskScheduler) ScheduleMultiScopedTxnTaskWithObserver(
	cfg *Config,
	taskType TaskType,
	scopes []common.ID,
	factory TxnTaskFactory,
	observers ...Observer) (task Task, err error) {
	task = NewScheduledTxnTask(s.ctx, s.txnMgr, taskType, cfg, scopes, factory)
	for _, observer := range observers {
		task.AddObserver(observer)
	}
	err = s.Schedule(task)
	return
}

func (s *taskScheduler) CheckAsyncScopes(scopes []common.ID) (err error) {
	dispatcher := s.typedHandler.handlers[DataCompactionTask].(*ScopeCheckTaskHandler)
	dispatcher.Lock()
	defer dispatcher.Unlock()
	return dispatcher.checkConflictLocked(scopes)
}

func (s *taskScheduler) ScheduleMultiScopedFn(
	ctx *Config,
	taskType TaskType,
	scopes []common.ID,
	fn FuncT) (task Task, err error) {
	task = NewMultiScopedFnTask(ctx, taskType, scopes, fn)
	err = s.Schedule(task)
	return
}

func (s *taskScheduler) GetPenddingLSNCnt() uint64 {
	return s.wal.GetPenddingCnt()
}

func (s *taskScheduler) GetCheckpointedLSN() uint64 {
	return s.wal.GetCheckpointed()
}

func (s *taskScheduler) ScheduleFn(ctx *Config, taskType TaskType, fn func() error) (task Task, err error) {
	task = NewFnTask(ctx, taskType, fn)
	err = s.Schedule(task)
	return
}

func (s *taskScheduler) ScheduleScopedFn(ctx *Config, taskType TaskType, scope *common.ID, fn func() error) (task Task, err error) {
	task = NewScopedFnTask(ctx, taskType, scope, fn)
	err = s.Schedule(task)
	return
}

func (s *taskScheduler) Schedule(task Task) (err error) {
	taskType := task.Type()
	// if taskType == DataCompactionTask || taskType == GCTask {
	if taskType == DataCompactionTask || taskType == FlushTableTailTask {
		dispatcher := s.typedHandler.handlers[DataCompactionTask].(*ScopeCheckTaskHandler)
		return dispatcher.TryDispatch(task)
	}
	s.typedHandler.Enqueue(task)
	return nil
}
