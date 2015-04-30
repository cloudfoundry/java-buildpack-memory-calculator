// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2015 the original author or authors.
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

import (
	"fmt"
	"math"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory/switches"
)

type Allocator interface {
	Balance(memLimit MemSize)        // Balance allocations to buckets within MemSize memory limit
	SetLowerBounds()                 // Allocate without a memory limit
	Switches(switches.Funs) []string // Generate memory switches from current allocations
}

type allocator struct {
	buckets map[string]Bucket
}

func NewAllocator(sizes map[string]string, heuristics map[string]float64) (*allocator, error) {
	if buckets, err := createMemoryBuckets(sizes, heuristics); err != nil {
		return nil, fmt.Errorf("allocator not created: %s", err)
	} else {
		return &allocator{
			buckets: buckets,
		}, nil
	}
}

func (a *allocator) Balance(memLimit MemSize) {
	// Adjust stack bucket, if it exists
	stackBucket, estNumThreads := a.normaliseStack(memLimit)

	// balance buckets
	a.balance()

	// Validate result and issue warnings?

	// Re-adjust stack bucket, if it exists
	a.unnormaliseStack(stackBucket, estNumThreads)
}

func (a *allocator) SetLowerBounds() {
	for _, b := range a.buckets {
		b.SetSize(b.Range().Floor())
	}
}

func (a *allocator) Switches(sfs switches.Funs) []string {
	var strs = make([]string, 0, 10)
	for s, b := range a.buckets {
		strs = append(strs, sfs.Apply(s, b.GetSize().String())...)
	}
	return strs
}

func createMemoryBuckets(sizes map[string]string, heuristics map[string]float64) (map[string]Bucket, error) {
	buckets := map[string]Bucket{}
	for name, weight := range heuristics {
		size, ok := sizes[name]
		if !ok {
			size = ".."
		}
		if aRange, err := NewRangeFromString(size); err == nil {
			if buckets[name], err = NewBucket(name, weight, aRange); err != nil {
				return nil, fmt.Errorf("memory type '%s' cannot be allocated: %s", name, err)
			}
		} else {
			return nil, fmt.Errorf("memory type '%s' cannot be allocated: %s", name, err)
		}
	}
	return buckets, nil
}

func (a *allocator) totalWeight() float64 {
	var w float64
	for _, b := range a.buckets {
		w = w + b.Weight()
	}
	return w
}

func (a *allocator) normaliseStack(memLimit MemSize) (Bucket, float64) {
	if sb, ok := a.buckets["stack"]; ok {
		stackMem := a.weightedProportion(memLimit, sb)
		estNum := math.Max(1.0, stackMem/float64(sb.DefaultSize()))
		nsb, _ := NewBucket("normalised stack", sb.Weight(), sb.Range().Scale(estNum))

		a.buckets["stack"] = nsb
		return sb, estNum
	}
	return nil, 0.0
}

func (a *allocator) weightedProportion(memLimit MemSize, b Bucket) float64 {
	return float64(memLimit) * b.Weight() / a.totalWeight()
}

func (a *allocator) unnormaliseStack(sb Bucket, estNum float64) {
	if sb == nil {
		return
	}
	newSize := (*a.buckets["stack"].GetSize()).Scale(1.0 / estNum)
	sb.SetSize(newSize)
	a.buckets["stack"] = sb
}

func (a *allocator) balance() {

}
