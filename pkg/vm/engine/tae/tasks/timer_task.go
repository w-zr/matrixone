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
	"time"
)

type lambdaHandle struct {
	onExec func()
	onStop func()
}

func (h *lambdaHandle) OnExec() {
	if h.onExec != nil {
		h.onExec()
	}
}

func (h *lambdaHandle) OnStopped() {
	if h.onStop != nil {
		h.onStop()
	}
}

type TimerTask struct {
	handle   taskHandle
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewTimerTaskWithFunc(interval time.Duration, onExec, onStop func()) *TimerTask {
	return NewTimerTask(interval, &lambdaHandle{onExec: onExec, onStop: onStop})
}

func NewTimerTask(interval time.Duration, handle taskHandle) *TimerTask {
	c := &TimerTask{
		interval: interval,
		handle:   handle,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

func (c *TimerTask) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
			}
			c.handle.OnExec()
		}
	}()
}

func (c *TimerTask) Stop() {
	c.cancel()
	c.wg.Wait()
	c.handle.OnStopped()
}
