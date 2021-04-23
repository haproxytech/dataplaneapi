// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package dataplaneapi

import (
	"context"
	"sync"
)

var (
	ContextHandler *CtxHandler
	once           sync.Once
)

type CtxHandler struct {
	ctx    context.Context
	rwx    sync.RWMutex
	cancel context.CancelFunc
}

func (ch *CtxHandler) Init() {
	once.Do(func() {
		ContextHandler = &CtxHandler{}
		ch = ContextHandler
	})
	ch.rwx.Lock()
	ch.ctx, ch.cancel = context.WithCancel(context.Background())
	ch.rwx.Unlock()
}

func (ch *CtxHandler) Cancel() {
	ch.rwx.RLock()
	ch.cancel()
	ch.rwx.RUnlock()
}

func (ch *CtxHandler) Context() context.Context {
	ch.rwx.RLock()
	defer ch.rwx.RUnlock()
	return ch.ctx
}
