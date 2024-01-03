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

package calculator

import (
	"fmt"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/v4/flags"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/v4/memory"
)

type Calculator struct {
	HeadRoom         *flags.HeadRoom
	JvmOptions       *flags.JVMOptions
	LoadedClassCount *flags.LoadedClassCount
	ThreadCount      *flags.ThreadCount
	TotalMemory      *flags.TotalMemory
}

func (c Calculator) Calculate() ([]fmt.Stringer, error) {
	var options []fmt.Stringer

	j := c.JvmOptions
	if j == nil {
		j = &flags.JVMOptions{}
	}

	headRoom := c.headRoom()

	directMemory := j.MaxDirectMemory
	if directMemory != nil {
		options = append(options, *directMemory)
	}

	metaspace := j.MaxMetaspace
	if metaspace == nil {
		m := c.metaspace()
		metaspace = &m
		options = append(options, *metaspace)
	}

	reservedCodeCache := j.ReservedCodeCache
	if reservedCodeCache == nil {
		r := memory.DefaultReservedCodeCache
		reservedCodeCache = &r
		options = append(options, *reservedCodeCache)
	}

	stack := j.Stack
	if stack == nil {
		s := memory.DefaultStack
		stack = &s
		options = append(options, *stack)
	}

	overhead := c.overhead(headRoom, metaspace, reservedCodeCache, stack)
	available := memory.Size(*c.TotalMemory)

	if overhead > available {
		return nil, fmt.Errorf("required memory %s is greater than %s available for allocation: %s, %s, %s x %d threads",
			overhead, available, metaspace, reservedCodeCache, stack, *c.ThreadCount)
	}

	heap := j.MaxHeap
	if heap == nil {
		h := c.heap(overhead)
		heap = &h
		options = append(options, *heap)
	}

	if overhead+memory.Size(*heap) > available {
		return nil, fmt.Errorf("required memory %s is greater than %s available for allocation: %s, %s, %s, %s x %d threads",
			overhead+memory.Size(*heap), available, heap, metaspace, reservedCodeCache, stack, *c.ThreadCount)
	}

	return options, nil
}

func (c Calculator) headRoom() memory.Size {
	return memory.Size(float64(*c.TotalMemory) * (float64(*c.HeadRoom) / 100))
}

func (c Calculator) heap(overhead memory.Size) memory.MaxHeap {
	return memory.MaxHeap(memory.Size(*c.TotalMemory) - overhead)
}

func (c Calculator) metaspace() memory.MaxMetaspace {
	return memory.MaxMetaspace((*c.LoadedClassCount * 5800) + 14000000)
}

func (c Calculator) overhead(headRoom memory.Size, metaspace *memory.MaxMetaspace, reservedCodeCache *memory.ReservedCodeCache, stack *memory.Stack) memory.Size {
	return headRoom +
		memory.Size(*metaspace) +
		memory.Size(*reservedCodeCache) +
		memory.Size(int64(*stack)*int64(*c.ThreadCount))
}
