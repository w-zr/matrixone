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
	"sync/atomic"
	"time"

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/logutil"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/common"
)

var WaitableCtx = &Config{Waitable: true}

type Config struct {
	ID       uint64
	DoneCB   opExecFunc
	Waitable bool
}

// func NewWaitableCtx() *Config {
// 	return &Config{Waitable: true}
// }

type BaseTask struct {
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

	parent   Task
	id       uint64
	taskType TaskType
	exec     func(Task) error
}

func NewBaseTask(parent Task, taskType TaskType, cfg *Config) *BaseTask {
	var id uint64
	if cfg != nil && cfg.ID != 0 {
		id = cfg.ID
	} else {
		id = NextTaskId()
	}
	task := &BaseTask{
		id:       id,
		taskType: taskType,
		parent:   parent,
	}
	var doneCB opExecFunc
	if cfg == nil {
		doneCB = task.onDone
	} else {
		if cfg.DoneCB == nil && !cfg.Waitable {
			doneCB = task.onDone
		}
	}
	if parent == nil {
		parent = task
	}
	task.impl = parent.(IOpInternal)
	task.doneCB = doneCB
	if doneCB == nil {
		task.errorC = make(chan error, 1)
	}
	return task
}

func (task *BaseTask) onDone(Task) {
	logutil.Debug("[Done]", common.OperationField(task.parent.Name()),
		common.DurationField(time.Duration(task.GetExecuteTime())),
		common.ErrorField(task.err))
}
func (task *BaseTask) Type() TaskType      { return task.taskType }
func (task *BaseTask) Cancel() (err error) { panic("todo") }
func (task *BaseTask) ID() uint64          { return task.id }
func (task *BaseTask) Execute() (err error) {
	if task.exec != nil {
		return task.exec(task)
	}
	logutil.Debugf("Execute Task Type=%d, ID=%d", task.taskType, task.id)
	return nil
}
func (task *BaseTask) Name() string {
	return fmt.Sprintf("Task[ID=%d][T=%s]", task.id, TaskName(task.taskType))
}

func (task *BaseTask) GetError() error {
	return task.err
}

func (task *BaseTask) SetError(err error) {
	task.endTime = time.Now()
	task.err = err
	if task.errorC != nil {
		task.errorC <- err
	} else if task.doneCB != nil {
		task.doneCB(task)
	} else {
		panic("logic error")
	}
	if task.observers != nil {
		for _, observer := range task.observers {
			observer.OnExecDone(task.impl)
		}
	}
}

func (task *BaseTask) Waitable() bool {
	return task.doneCB == nil
}

func (task *BaseTask) WaitDone(ctx context.Context) error {
	if task.waitedOnce.Load() {
		return moerr.NewTAEErrorNoCtx("wait done twice")
	}
	defer task.waitedOnce.Store(true)

	if task.errorC == nil {
		return moerr.NewTAEErrorNoCtx("wait done without error channel")
	}
	select {
	case err := <-task.errorC:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (task *BaseTask) PreExecute() error {
	return nil
}

func (task *BaseTask) PostExecute() error {
	return nil
}

func (task *BaseTask) OnExec(ctx context.Context) error {
	task.startTime = time.Now()
	err := task.impl.PreExecute()
	if err != nil {
		return err
	}
	err = task.impl.Execute(ctx)
	if err != nil {
		return err
	}
	err = task.impl.PostExecute()
	return err
}

func (task *BaseTask) GetCreateTime() time.Time {
	return task.createTime
}

func (task *BaseTask) GetStartTime() time.Time {
	return task.startTime
}

func (task *BaseTask) GetEndTime() time.Time {
	return task.endTime
}

func (task *BaseTask) GetExecuteTime() int64 {
	return task.endTime.Sub(task.startTime).Microseconds()
}

func (task *BaseTask) AddObserver(o Observer) {
	if task.observers == nil {
		task.observers = make([]Observer, 0)
	}
	task.observers = append(task.observers, o)
}
