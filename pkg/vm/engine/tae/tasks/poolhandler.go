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
	"sync"

	"github.com/matrixorigin/matrixone/pkg/logutil"
	"github.com/panjf2000/ants/v2"
)

var (
	poolHandlerName = "PoolHandler"
)

type PoolHandler struct {
	wg   sync.WaitGroup
	pool *ants.Pool

	opExec      opExecFunc
	taskHandler *baseTaskHandler
}

func NewPoolHandler(ctx context.Context, num int) *PoolHandler {
	pool, err := ants.NewPool(num)
	if err != nil {
		panic(err)
	}
	h := &PoolHandler{
		taskHandler: NewBaseTaskHandler(ctx, poolHandlerName),
		pool:        pool,
	}
	h.opExec = h.taskHandler.queue.execFunc
	h.taskHandler.queue.execFunc = h.doHandle
	return h
}

func (p *PoolHandler) Enqueue(task Task) {
	p.taskHandler.Enqueue(task)
}

func (p *PoolHandler) Start() {
	p.taskHandler.queue.Start()
}

func (h *PoolHandler) Close() error {
	h.pool.Release()
	h.taskHandler.Close()
	h.wg.Wait()
	return nil
}

func (h *PoolHandler) doHandle(task Task) {
	closure := func(o Task) func() {
		return func() {
			h.opExec(o)
			h.wg.Done()
		}
	}
	h.wg.Add(1)
	err := h.pool.Submit(closure(task))
	if err != nil {
		logutil.Warnf("%v", err)
		task.SetError(err)
		h.wg.Done()
	}
}
