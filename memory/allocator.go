// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2015-2017 the original author or authors.
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

package memory

import "fmt"

type Allocator interface {
	Calculate(loadedClasses int, stackThreads int, memLimit MemSize) error // Calculate vm options
	String() string                                                        // Get vm options (if calculation succeeded)
}

type allocator struct {
	vmOptions  VmOptions
	estimators map[MemoryType]func(int) MemSize
}

func NewAllocator(poolType string, vmOptions VmOptions) (*allocator, error) {
	if poolType == "metaspace" {
		return &allocator{
			vmOptions:  vmOptions,
			estimators: estimators,
		}, nil
	} else {
		return &allocator{
			vmOptions:  vmOptions,
			estimators: permgenEstimators,
		}, nil
	}
}

const (
	DEFAULT_MAX_DIRECT_MEMORY_SIZE int64 = 10 * 1024 * 1024
	DEFAULT_STACK_SIZE             int64 = 1024 * 1024
	DEFAULT_CODE_CACHE_JAVA_8      int64 = 240 * 1024 * 1024
)

var estimators = map[MemoryType]func(int) MemSize{
	MaxDirectMemorySize: func(loadedClasses int) MemSize {
		return NewMemSize(DEFAULT_MAX_DIRECT_MEMORY_SIZE)
	},
	MaxMetaspaceSize: func(loadedClasses int) MemSize {
		return NewMemSize(5400).Scale(float64(loadedClasses)).Add(NewMemSize(7000000))
	},
	ReservedCodeCacheSize: func(loadedClasses int) MemSize {
		return NewMemSize(DEFAULT_CODE_CACHE_JAVA_8)
	},
	CompressedClassSpaceSize: func(loadedClasses int) MemSize {
		return NewMemSize(700).Scale(float64(loadedClasses)).Add(NewMemSize(750000))
	},
}

const DEFAULT_CODE_CACHE_JAVA_7 int64 = 48 * 1024 * 1024

var permgenEstimators = map[MemoryType]func(int) MemSize{
	MaxDirectMemorySize: func(loadedClasses int) MemSize {
		return NewMemSize(DEFAULT_MAX_DIRECT_MEMORY_SIZE)
	},
	MaxPermSize: func(loadedClasses int) MemSize {
		return NewMemSize(6000).Scale(float64(loadedClasses)).Add(NewMemSize(7000000))
	},
	ReservedCodeCacheSize: func(loadedClasses int) MemSize {
		return NewMemSize(DEFAULT_CODE_CACHE_JAVA_7)
	},
}

func (a *allocator) Calculate(loadedClasses int, stackThreads int, memLimit MemSize) error {
	if memLimit.LessThan(MemSize(kILO)) {
		return fmt.Errorf("Too little memory to allocate: %s", memLimit)
	}

	for memoryType, estimator := range a.estimators {
		if !a.present(memoryType) {
			a.vmOptions.SetMemOpt(memoryType, estimator(loadedClasses))
		}
	}

	if !a.present(MaxHeapSize) {
		maxHeapSize, err := a.calculateMaxHeapSize(stackThreads, memLimit)
		if err != nil {
			return err
		}
		a.vmOptions.SetMemOpt(MaxHeapSize, maxHeapSize)
	}

	return nil
}

func (a *allocator) calculateMaxHeapSize(stackThreads int, memLimit MemSize) (MemSize, error) {
	var stackSize MemSize
	if a.present(StackSize) {
		stackSize = a.vmOptions.MemOpt(StackSize)
	} else {
		stackSize = NewMemSize(DEFAULT_STACK_SIZE)
	}

	allocatedMemory := stackSize.Scale(float64(stackThreads)).Add(a.estimatedMemory())

	maxHeapSize := memLimit.Subtract(allocatedMemory)
	if maxHeapSize.LessThan(MEMSIZE_ZERO) {
		ss := ""
		if !a.present(StackSize) {
			ss = fmt.Sprintf(" %s%s", switches[StackSize], stackSize.String())
		}
		return MEMSIZE_ZERO, fmt.Errorf("insufficient memory remaining for heap (memory limit %s < allocated memory %s): %s%s", memLimit, allocatedMemory, a.vmOptions.String(), ss)
	}

	return maxHeapSize, nil
}

func (a *allocator) estimatedMemory() MemSize {
	est := NewMemSize(0)
	for memoryType, _ := range a.estimators {
		est = est.Add(a.vmOptions.MemOpt(memoryType))
	}
	return est
}

func (a *allocator) present(memoryType MemoryType) bool {
	return !a.vmOptions.MemOpt(memoryType).Equals(MEMSIZE_ZERO)
}

func (a *allocator) String() string {
	return a.vmOptions.DeltaString()
}
