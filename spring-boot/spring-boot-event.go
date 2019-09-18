/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package SpringBoot

import (
	"sync"
	"github.com/didi/go-spring/spring-core"
)

type SpringApplicationContext interface {
	SpringCore.SpringContext

	Wait()
	SafeGoroutine(f func())
}

type SpringApplicationEvent interface {
	OnStartApplication(ctx SpringApplicationContext)
	OnStopApplication(ctx SpringApplicationContext)
}

type DefaultSpringApplicationContext struct {
	*SpringCore.DefaultSpringContext

	wg sync.WaitGroup
}

func NewSpringApplicationContext() *DefaultSpringApplicationContext {
	return &DefaultSpringApplicationContext{
		DefaultSpringContext: SpringCore.NewSpringContext(),
	}
}

func (ctx *DefaultSpringApplicationContext) SafeGoroutine(f func()) {
	go func() {

		defer func() {

		}()

		ctx.wg.Add(1)
		defer ctx.wg.Done()

		f()
	}()
}

func (ctx *DefaultSpringApplicationContext) Wait() {
	ctx.wg.Wait()
}
