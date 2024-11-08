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

	"github.com/matrixorigin/matrixone/pkg/common/moerr"
	"github.com/matrixorigin/matrixone/pkg/logutil"
	"github.com/matrixorigin/matrixone/pkg/vm/engine/tae/common"
)

type Cmd = uint8

const (
	QUIT Cmd = iota
)

const (
	CREATED int32 = iota
	RUNNING
	StoppingReceiver
	StoppingCMD
	STOPPED
)

const (
	QueueSize = 10000
)

type opExecFunc func(Task)

type stats struct {
	processed atomic.Uint64
	succeeded atomic.Uint64
	failed    atomic.Uint64
	avgTime   atomic.Int64
}

func (s *stats) addProcessed() {
	s.processed.Add(1)
}

func (s *stats) addSucceeded() {
	s.succeeded.Add(1)
}

func (s *stats) addFailed() {
	s.failed.Add(1)
}

func (s *stats) recordTime(t int64) {
	procced := s.processed.Load()
	avg := s.avgTime.Load()
	//TODO: avgTime is wrong
	s.avgTime.Store((avg*int64(procced-1) + t) / int64(procced))
}

func (s *stats) String() string {
	r := fmt.Sprintf("Total: %d, Succ: %d, Fail: %d, avgTime: %dus",
		s.processed.Load(),
		s.failed.Load(),
		s.avgTime.Load(),
		s.avgTime.Load())
	return r
}

type OpQueue struct {
	name       string
	taskC      chan Task
	cmdC       chan Cmd
	state      atomic.Int32
	pending    atomic.Int64
	closedC    chan struct{}
	stats      stats
	execFunc   opExecFunc
	cancelFunc opExecFunc
}

func NewOpQueue(ctx context.Context, name string, args ...int) *OpQueue {
	var l int
	if len(args) == 0 {
		l = QueueSize
	} else {
		l = args[0]
		if l < 0 {
			logutil.Warnf("Create OpQueue with negtive queue size %d", l)
			l = QueueSize
		}
	}
	if name == "" {
		name = fmt.Sprintf("[queue-%d]", common.NextGlobalSeqNum())
	}
	queue := &OpQueue{
		name:    name,
		taskC:   make(chan Task, l),
		cmdC:    make(chan Cmd, l),
		closedC: make(chan struct{}),
	}
	queue.state.Store(CREATED)

	queue.execFunc = func(task Task) {
		err := task.OnExec(ctx)
		queue.stats.addProcessed()
		if err != nil {
			queue.stats.addFailed()
		} else {
			queue.stats.addSucceeded()
		}
		task.SetError(err)
		queue.stats.recordTime(task.GetExecuteTime())
	}

	queue.cancelFunc = func(task Task) {
		task.SetError(moerr.NewInternalErrorNoCtx("op cancelled"))
	}
	return queue
}

func (w *OpQueue) Start() {
	logutil.Debugf("%s Started", w.name)
	if w.state.Load() != CREATED {
		if w.state.Load() == RUNNING {
			logutil.Warnf("op queue has started")
			return
		} else {
			panic(fmt.Sprintf("logic error: %v", w.state.Load()))
		}
	}
	w.state.Store(RUNNING)
	go func() {
		for {
			state := w.state.Load()
			if state == STOPPED {
				break
			}
			select {
			case op := <-w.taskC:
				w.execFunc(op)
				// if state == RUNNING {
				// 	w.execFunc(op)
				// } else {
				// 	w.CancelFunc(op)
				// }
				w.pending.Add(-1)
			case cmd := <-w.cmdC:
				w.onCmd(cmd)
			}
		}
	}()
}

func (w *OpQueue) Stop() {
	w.StopReceiver()
	w.WaitStop()
	logutil.Debugf("%s Stopped", w.name)
}

func (w *OpQueue) StopReceiver() {
	state := w.state.Load()
	if state >= StoppingReceiver {
		return
	}
	w.state.CompareAndSwap(state, StoppingReceiver)
}

func (w *OpQueue) WaitStop() {
	state := w.state.Load()
	if state <= RUNNING {
		panic("logic error")
	}
	if state == STOPPED {
		return
	}
	if w.state.CompareAndSwap(StoppingReceiver, StoppingCMD) {
		pending := w.pending.Load()
		for {
			if pending == 0 {
				break
			}
			pending = w.pending.Load()
		}
		w.cmdC <- QUIT
	}
	<-w.closedC
}

func (w *OpQueue) Enqueue(task Task) bool {
	state := w.state.Load()
	if state != RUNNING {
		return false
	}
	w.pending.Add(1)
	if w.state.Load() != RUNNING {
		w.pending.Add(-1)
		return false
	}
	w.taskC <- task
	return true
}

func (w *OpQueue) onCmd(cmd Cmd) {
	switch cmd {
	case QUIT:
		// log.Infof("Quit OpQueue")
		close(w.cmdC)
		close(w.taskC)
		if !w.state.CompareAndSwap(StoppingCMD, STOPPED) {
			panic("logic error")
		}
		w.closedC <- struct{}{}
	default:
		panic(fmt.Sprintf("Unsupported cmd %d", cmd))
	}
}

func (w *OpQueue) StatsString() string {
	return fmt.Sprintf("| stats | %s | w | %s", w.stats.String(), w.name)
}
