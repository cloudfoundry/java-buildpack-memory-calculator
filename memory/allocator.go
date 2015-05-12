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
	Balance(memLimit MemSize) error  // Balance allocations to buckets within memory limit
	Switches(switches.Funs) []string // Get selected memory switches from current allocations
	GetWarnings() []string           // Get warnings (if balancing succeeded)
}

type allocator struct {
	originalSizes map[string]Range  // unmodified after creation
	buckets       map[string]Bucket // named buckets for allocation
	warnings      []string          // warnings if allocation found issues
}

func NewAllocator(sizes map[string]Range, heuristics map[string]float64) (*allocator, error) {
	if buckets, err := createMemoryBuckets(sizes, heuristics); err != nil {
		return nil, fmt.Errorf("allocator not created: %s", err)
	} else {
		return &allocator{
			originalSizes: sizes,
			buckets:       buckets,
		}, nil
	}
}

const (
	NATIVE_MEMORY_WARNING_FACTOR float64 = 3.0
	TOTAL_MEMORY_WARNING_FACTOR  float64 = 0.8
	CLOSE_TO_DEFAULT_FACTOR      float64 = 0.1
)

// Balance memory between buckets, adjusting stack units, observing
// constraints, and detecting memory wastage and default proximity.
func (a *allocator) Balance(memLimit MemSize) error {
	if memLimit.LessThan(MemSize(kILO)) {
		return fmt.Errorf("Too little memory to allocate: %s", memLimit)
	}

	// adjust stack bucket, if it exists
	stackBucket, estNumThreads := a.normaliseStack(memLimit)

	// distribute memory among the buckets
	if berr := a.balance(memLimit); berr != nil {
		return fmt.Errorf("Memory allocation failed for configuration: %v, : %s", getSizes(a.originalSizes), berr)
	}

	// validate result and gather warnings
	a.validateAllocation(memLimit)

	// reset stack bucket, if it exists
	a.unnormaliseStack(stackBucket, estNumThreads)

	return nil
}

func (a *allocator) Switches(sfs switches.Funs) []string {
	var strs = make([]string, 0, 10)
	for s, b := range a.buckets {
		strs = append(strs, sfs.Apply(s, b.GetSize().String())...)
	}
	return strs
}

func (a *allocator) GetWarnings() []string {
	return a.warnings
}

// getSizes returns a slice of memory type range strings
func getSizes(ss map[string]Range) []string {
	result := []string{}
	for n, s := range ss {
		result = append(result, n+":"+s.String())
	}
	return result
}

func createMemoryBuckets(sizes map[string]Range, heuristics map[string]float64) (map[string]Bucket, error) {
	buckets := map[string]Bucket{}
	for name, weight := range heuristics {
		aRange, ok := sizes[name]
		if !ok {
			aRange, _ = NewUnboundedRange(MEMSIZE_ZERO)
		}
		var err error
		if buckets[name], err = NewBucket(name, weight, aRange); err != nil {
			return nil, fmt.Errorf("memory type '%s' cannot be allocated: %s", name, err)
		}
	}
	return buckets, nil
}

func totalWeight(bs map[string]Bucket) float64 {
	var w float64
	for _, b := range bs {
		w = w + b.Weight()
	}
	return w
}

// Replace stack bucket to make it represent total memory for stacks temporarily
func (a *allocator) normaliseStack(memLimit MemSize) (originalStackBucket Bucket, estNumThreads float64) {
	if sb, ok := a.buckets["stack"]; ok {
		stackMem := weightedSize(totalWeight(a.buckets), memLimit, sb)
		estNumThreads = math.Max(1.0, stackMem/float64(sb.DefaultSize()))
		nsb, _ := NewBucket("normalised stack", sb.Weight(), sb.Range().Scale(estNumThreads))

		a.buckets["stack"] = nsb
		return sb, estNumThreads
	}
	return nil, 0.0
}

func weightedSize(totWeight float64, memLimit MemSize, b Bucket) float64 {
	return (float64(memLimit) * b.Weight()) / totWeight
}

// Replace stack bucket, and set size per thread
func (a *allocator) unnormaliseStack(sb Bucket, estNum float64) {
	if sb == nil {
		return
	}
	newSize := (*a.buckets["stack"].GetSize()).Scale(1.0 / estNum)
	sb.SetSize(newSize)
	a.buckets["stack"] = sb
}

// Balance memory between buckets, observing constraints.
func (a *allocator) balance(memLimit MemSize) error {
	remaining := copyBucketMap(a.buckets)
	removed := true

	for removed && len(remaining) != 0 {
		var err error
		memLimit, removed, err = balanceOrRemove(remaining, memLimit)
		if err != nil {
			return err
		}
	}

	// check for zero allocations
	for n, b := range a.buckets {
		if b.GetSize().String() == "0" {
			return fmt.Errorf("Cannot allocate memory to '%n' type", n)
		}
	}

	return nil
}

func (a *allocator) validateAllocation(memLimit MemSize) {
	// memory_wastage_warning(buckets)
	a.warnings = append(a.warnings, memoryWastageWarnings(a.buckets, memLimit)...)
	// close_to_default_warnings(buckets)
	a.warnings = append(a.warnings, closeToDefaultWarnings(a.buckets, a.originalSizes, memLimit)...)
}

func memoryWastageWarnings(bs map[string]Bucket, memLimit MemSize) []string {
	warnings := []string{}
	if nb, ok := bs["native"]; ok {
		warnings = append(warnings, nativeBucketWarning(nb, totalWeight(bs), memLimit)...)
	}

	totalSize := MEMSIZE_ZERO
	for _, b := range bs {
		totalSize = totalSize.Add(*b.GetSize())
	}
	if totalSize.LessThan(memLimit.Scale(TOTAL_MEMORY_WARNING_FACTOR)) {
		warnings = append(warnings,
			fmt.Sprintf(
				"The allocated Java memory sizes total %s which is less than %g of "+
					"the available memory, so configured Java memory sizes may be too small or available memory may be too large",
				totalSize.String(), TOTAL_MEMORY_WARNING_FACTOR))
	}
	return warnings
}

func nativeBucketWarning(nativeBucket Bucket, totWeight float64, memLimit MemSize) []string {
	if nativeBucket.Range().Floor() == MEMSIZE_ZERO {
		floatSize := NATIVE_MEMORY_WARNING_FACTOR * weightedSize(totWeight, memLimit, nativeBucket)
		if MemSize(math.Floor(0.5 + floatSize)).LessThan(*nativeBucket.GetSize()) {
			return []string{fmt.Sprintf(
				"There is more than %g times more spare native memory than the default "+
					"so configured Java memory may be too small or available memory may be too large",
				NATIVE_MEMORY_WARNING_FACTOR)}
		}
	}
	return []string{}
}

func closeToDefaultWarnings(bs map[string]Bucket, sizes map[string]Range, memLimit MemSize) []string {
	totWeight := totalWeight(bs)
	warnings := []string{}
	for name, b := range bs {
		if _, ok := sizes[name]; ok && name != "stack" && b.Range().Degenerate() {
			floatDefSize := weightedSize(totWeight, memLimit, b)
			floatActualSize := float64(*b.GetSize())
			var factor float64
			if floatDefSize > 0.0 {
				factor = math.Abs((floatActualSize - floatDefSize) / floatDefSize)
			}
			if (floatDefSize == 0.0 && floatActualSize == 0.0) || (factor < CLOSE_TO_DEFAULT_FACTOR) {
				warnings = append(warnings,
					fmt.Sprintf(
						"The specified value %s for memory type %s is close to the computed value %s. "+
							"Consider taking the default.",
						b.GetSize(), name, MemSize(math.Floor(0.5+floatDefSize))))
			}
		}
	}
	return warnings
}

// Balance the allocation of memLeft memory between remaining buckets, and see
// if any buckets cannot be so allocated. Buckets that are constrained are
// then allocated the constrained amount, and removed from the remaining map.
//
// This function modifies the caller's bucket map.
//
// Returns the remaining memory not so allocated, a bool indicating if some were
// removed from consideration, and a possible memory exceeded error.
func balanceOrRemove(remaining map[string]Bucket, memLeft MemSize) (MemSize, bool, error) {
	bucketsRemoved := []string{}
	remainingWeight := totalWeight(remaining)

	allocatedInThisPass := MEMSIZE_ZERO
	for name, b := range remaining {

		// round to nearest, so as not to lose a byte inadvertently.
		size := MemSize(math.Floor(0.5 + weightedSize(remainingWeight, memLeft, b)))

		if b.Range().Contains(size) {
			// speculatively set the size, in case this pass doesn't remove any buckets
			b.SetSize(size)
		} else {
			newSize := b.Range().Constrain(size)
			b.SetSize(newSize)
			allocatedInThisPass = allocatedInThisPass.Add(newSize)
			bucketsRemoved = append(bucketsRemoved, name)
		}
	}

	memLeft = memLeft.Subtract(allocatedInThisPass)
	if memLeft.LessThan(MEMSIZE_ZERO) {
		return 0, false, fmt.Errorf("memory exceeded")
	}

	// do deletes afterwards in case the for .. range is disturbed
	for _, name := range bucketsRemoved {
		delete(remaining, name)
	}

	return memLeft, len(bucketsRemoved) != 0, nil
}

// Shallow copy of bucket map.
func copyBucketMap(bs map[string]Bucket) map[string]Bucket {
	cpy := map[string]Bucket{}
	for k, b := range bs {
		cpy[k] = b
	}
	return cpy
}
