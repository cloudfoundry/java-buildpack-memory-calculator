/*
 * Copyright 2015-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package flags

import (
	"strings"

	"github.com/instana/java-buildpack-memory-calculator/v4/memory"
)

var DefaultJVMOptions = JVMOptions{}

const FlagJVMOptions = "jvm-options"

type JVMOptions struct {
	MaxDirectMemory   *memory.MaxDirectMemory
	MaxHeap           *memory.MaxHeap
	MaxMetaspace      *memory.MaxMetaspace
	ReservedCodeCache *memory.ReservedCodeCache
	Stack             *memory.Stack
}

func (j *JVMOptions) Set(s string) error {
	for _, c := range strings.Split(s, " ") {
		if memory.IsMaxDirectMemory(c) {
			m, err := memory.ParseMaxDirectMemory(c)
			if err != nil {
				return err
			}

			j.MaxDirectMemory = &m
		} else if memory.IsMaxHeap(c) {
			m, err := memory.ParseMaxHeap(c)
			if err != nil {
				return err
			}

			j.MaxHeap = &m
		} else if memory.IsMaxMetaspace(c) {
			m, err := memory.ParseMaxMetaspace(c)
			if err != nil {
				return err
			}

			j.MaxMetaspace = &m
		} else if memory.IsReservedCodeCache(c) {
			r, err := memory.ParseReservedCodeCache(c)
			if err != nil {
				return err
			}

			j.ReservedCodeCache = &r
		} else if memory.IsStack(c) {
			s, err := memory.ParseStack(c)
			if err != nil {
				return err
			}

			j.Stack = &s
		}
	}

	return nil
}

func (j *JVMOptions) String() string {
	var values []string

	if j.MaxDirectMemory != nil {
		values = append(values, j.MaxDirectMemory.String())
	}

	if j.MaxHeap != nil {
		values = append(values, j.MaxHeap.String())
	}

	if j.MaxMetaspace != nil {
		values = append(values, j.MaxMetaspace.String())
	}

	if j.ReservedCodeCache != nil {
		values = append(values, j.ReservedCodeCache.String())
	}

	if j.Stack != nil {
		values = append(values, j.Stack.String())
	}

	return strings.Join(values, " ")
}

func (j *JVMOptions) Type() string {
	return "string"
}

func (j *JVMOptions) Validate() error {
	return nil
}
