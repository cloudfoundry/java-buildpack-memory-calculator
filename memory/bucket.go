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

// This file defines and implements the Bucket interface type.
// A Bucket is used to calculate default sizes for various type of memory.

type Bucket interface {
	Name() string // Name of bucket

	GetSize() *MemSize // Size of bucket, if set; nil otherwise
	SetSize(MemSize)   // cannot unset the size, once set

	Range() Range // Permissible range of sizes for this bucket.
	SetRange(Range)

	Weight() float64 // Proportion of total memory this bucket is allowed to consume by default.

	DefaultSize() MemSize // only supported by 'stack' buckets
}

type bucket struct {
	name   string
	size   *MemSize
	srange Range
	weight float64
}

// Returns the name of the bucket
func (b *bucket) Name() string {
	return b.name
}

// Returns a pointer to MemSize because this value can be unset.
func (b *bucket) GetSize() *MemSize {
	return b.size
}

// Generates a new pointer to MemSize internally.
func (b *bucket) SetSize(size MemSize) {
	tmpsize := size
	b.size = &tmpsize
}

func (b *bucket) Range() Range {
	return b.srange
}

func (b *bucket) SetRange(srange Range) {
	b.srange = srange
}

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
	if weight < 0.0 || weight > 1.0 {
		return nil, fmt.Errorf("Weight (%g) for bucket %s must be <=1 and >=0.", weight, name)
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

// The default stacksize: Floor() of the range, or the JRE standard default if
// the Floor is 0. Zero if this is not a bucket named 'stack'.
func (b *bucket) DefaultSize() MemSize {
	if b.name != "stack" {
		return MS_ZERO
	}
	floor := b.srange.Floor()
	if floor.Bytes() == 0 {
		return jre_DEFAULT_STACK_SIZE
	}
	return floor
}
