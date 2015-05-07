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
	"strings"
)

// A Bucket is used to calculate default sizes for various type of memory.
type Bucket interface {
	Name() string         // Name of bucket
	GetSize() *MemSize    // Size of bucket, if set; nil if not set
	SetSize(MemSize)      // (cannot unset the size, once set)
	Range() Range         // Permissible range of sizes for this bucket.
	Weight() float64      // Proportion of total memory this bucket is allowed to consume by default.
	DefaultSize() MemSize // Default size for 'stack' buckets
	String() string
}

type bucket struct {
	name   string
	size   *MemSize
	srange Range
	weight float64
}

// Name returns the (internal) bucket name
func (b *bucket) Name() string {
	return b.name
}

// GetSize returns a pointer to the size of the bucket because this value may not be set.
func (b *bucket) GetSize() *MemSize {
	return b.size
}

// SetSize sets the size of the bucket
func (b *bucket) SetSize(size MemSize) {
	tmpsize := size
	b.size = &tmpsize
}

// Range returns the bucket range constraint
func (b *bucket) Range() Range {
	return b.srange
}

// Weight returns the (relative) weight for this bucket allocation
func (b *bucket) Weight() float64 {
	return b.weight
}

func NewBucket(name string, weight float64, srange Range) (Bucket, error) {
	return newBucket(name, weight, srange)
}

func newBucket(name string, weight float64, srange Range) (*bucket, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("Bucket name must not be blank.")
	}
	if weight <= 0.0 {
		return nil, fmt.Errorf("Weight (%g) for bucket %s must be positive.", weight, name)
	}
	return &bucket{
		name:   name,
		srange: srange,
		weight: weight,
	}, nil
}

var (
	jre_DEFAULT_STACK_SIZE = NewMemSize(mEGA)
)

// The default stack size: Floor() of the Range, or the JRE standard default if
// the floor is 0. (MEMSIZE_ZERO if this is not a bucket named 'stack'.)
func (b *bucket) DefaultSize() MemSize {
	if b.name != "stack" {
		return MEMSIZE_ZERO
	}
	floor := b.srange.Floor()
	if floor.Bytes() == 0 {
		return jre_DEFAULT_STACK_SIZE
	}
	return floor
}

// String representation of a bucket, used for testing.
func (b *bucket) String() string {
	return fmt.Sprintf("Bucket{name: %s, size: %s, range: %s, weight: %g}", b.name, b.size, b.srange, b.weight)
}
