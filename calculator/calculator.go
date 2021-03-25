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
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/instana/java-buildpack-memory-calculator/v4/flags"
	"github.com/instana/java-buildpack-memory-calculator/v4/memory"
)

type Calculator struct {
	HeadRoom                 *flags.HeadRoom
	JvmOptions               *flags.JVMOptions
	LoadedClassCount         *flags.LoadedClassCount
	ThreadCount              *flags.ThreadCount
	DirectMemoryToHeapRatio  *flags.DirectMemoryToHeapRatio
	HeapYoungGenerationRatio *flags.HeapYoungGenerationRatio
	TotalMemory              *flags.TotalMemory
	DetectMemoryLimits       *flags.DetectMemoryLimits
}

const (
	LinuxMemoryLimitsPath        = "/sys/fs/cgroup/memory/memory.limit_in_bytes"
	LinuxUnspecifiedMemoryLimits = "9223372036854771712"
)

func (c Calculator) Calculate() ([]fmt.Stringer, error) {
	var options []fmt.Stringer

	j := c.JvmOptions
	if j == nil {
		j = &flags.JVMOptions{}
	}

	var available memory.Size

	if *c.TotalMemory == 0 && *c.DetectMemoryLimits == false {
		return nil, fmt.Errorf("neither '--%s' nor '--%s' are specified; cannot perform memory calculations without the total memory available",
			flags.FlagDetectMemoryLimits, flags.FlagTotalMemory)
	}

	if *c.DetectMemoryLimits == true {
		if *c.TotalMemory != 0 {
			return nil, fmt.Errorf("both '--%s' nor '--%s' are specified", flags.FlagDetectMemoryLimits, flags.FlagTotalMemory)
		}

		var err error
		available, err = c.detectMemoryLimits()
		if err != nil {
			return nil, fmt.Errorf("retrieving the memory limits automatically failed: %v", err)
		}
	} else {
		available = memory.Size(*c.TotalMemory)
	}

	headRoom := c.headRoom(available)

	directMemory := j.MaxDirectMemory
	directMemoryToHeapRatio := c.DirectMemoryToHeapRatio

	useDirectMemoryToHeapRatio := false

	if directMemory == nil && directMemoryToHeapRatio == nil {
		d := memory.DefaultMaxDirectMemory
		directMemory = &d
		options = append(options, *directMemory)
	} else if directMemoryToHeapRatio != nil {
		useDirectMemoryToHeapRatio = true
	}

	heapYoungGenerationRatio := c.HeapYoungGenerationRatio

	metaspace := j.MaxMetaspace
	if metaspace == nil {
		if *c.LoadedClassCount < 1 {
			return nil, fmt.Errorf("neither the '%s' argument nor the '%s' JVM option are specified; cannot calculate the Metaspace sizing", flags.FlagLoadedClassCount, "-XX:MaxMetaspaceSize")
		}

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

	overhead := c.overhead(headRoom, directMemory, metaspace, reservedCodeCache, stack)

	if overhead > available {
		return nil, fmt.Errorf("required memory %s is greater than %s available for allocation: %s, %s, %s, %s x %d threads",
			overhead, available, directMemory, metaspace, reservedCodeCache, stack, *c.ThreadCount)
	}

	var dynamicallyAllocatedMemory memory.Size

	heap := j.MaxHeap
	if heap != nil {
		dynamicallyAllocatedMemory = memory.Size(*heap)
	} else if useDirectMemoryToHeapRatio {
		// Split available memory between direct memory and heap
		availableMemory := available - overhead

		directMemorySize := (int64)(float64(availableMemory) * float64(*c.DirectMemoryToHeapRatio))
		heapMemorySize := availableMemory - memory.Size(directMemorySize)

		m := memory.MaxDirectMemory(memory.Size(directMemorySize))
		directMemory = &m

		h := memory.MaxHeap(memory.Size(heapMemorySize))
		heap = &h

		options = append(options, *heap)
		options = append(options, *directMemory)

		dynamicallyAllocatedMemory = memory.Size(directMemorySize) + memory.Size(heapMemorySize)
	} else {
		// Give all the available memory to the heap
		h := c.heap(available, overhead)
		heap = &h
		options = append(options, *heap)

		dynamicallyAllocatedMemory = memory.Size(h)
	}

	heapYoungGeneration := j.MaxHeapYoungGeneration
	if heapYoungGeneration == nil && heapYoungGenerationRatio != nil {
		youngGenerationSize := float32(*heap) * float32(*heapYoungGenerationRatio)
		y := memory.MaxHeapYoungGeneration(memory.Size(youngGenerationSize))
		youngGeneration := &y

		options = append(options, *youngGeneration)
	}

	if overhead+dynamicallyAllocatedMemory > available {
		return nil, fmt.Errorf("required memory %s is greater than %s available for allocation: %s, %s, %s, %s, %s x %d threads",
			overhead+dynamicallyAllocatedMemory, available, directMemory, heap, metaspace, reservedCodeCache, stack, *c.ThreadCount)
	}

	return options, nil
}

func (c Calculator) headRoom(available memory.Size) memory.Size {
	return memory.Size(float64(available) * (float64(*c.HeadRoom) / 100))
}

func (c Calculator) heap(available memory.Size, overhead memory.Size) memory.MaxHeap {
	return memory.MaxHeap(memory.Size(available) - overhead)
}

func (c Calculator) metaspace() memory.MaxMetaspace {
	return memory.MaxMetaspace((*c.LoadedClassCount * 5800) + 14000000)
}

func (c Calculator) overhead(headRoom memory.Size, directMemory *memory.MaxDirectMemory, metaspace *memory.MaxMetaspace, reservedCodeCache *memory.ReservedCodeCache, stack *memory.Stack) memory.Size {
	overhead := headRoom +
		memory.Size(*metaspace) +
		memory.Size(*reservedCodeCache) +
		memory.Size(int64(*stack)*int64(*c.ThreadCount))

	if directMemory != nil {
		overhead = overhead + memory.Size(*directMemory)
	}

	return overhead
}

func (c Calculator) detectMemoryLimits() (memory.Size, error) {
	if runtime.GOOS != "linux" {
		return 0, fmt.Errorf("the '--%s' option is supported only on Linux", flags.FlagDetectMemoryLimits)
	}

	bs, err := ioutil.ReadFile(LinuxMemoryLimitsPath)
	if err != nil {
		return 0, fmt.Errorf("cannot detect memory limit automatically from '%s': %v", LinuxMemoryLimitsPath, err)
	}

	v := strings.TrimSpace(string(bs))
	if v == LinuxUnspecifiedMemoryLimits {
		return 0, fmt.Errorf("no memory limits specified, '%s' contains '%s'", LinuxMemoryLimitsPath, LinuxUnspecifiedMemoryLimits)
	}

	s, err := memory.ParseSize(v + "b")
	if err != nil {
		return 0, fmt.Errorf("cannot parse as memory the '%s' output of '%s': %v", v, LinuxMemoryLimitsPath, err)
	}

	return s, nil
}
