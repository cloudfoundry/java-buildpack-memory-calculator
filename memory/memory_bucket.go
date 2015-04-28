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
	GetKSize() *int64 // Size of bucket in kilobytes, if set; nil otherwise
	SetKSize(int64)

	GetRange() Range // Permissible range of sizes for this bucket.
	SetRange(Range)

	GetWeight() float64 // Proportion of total memory this bucket is allowed to consume by default.
}

type bucket struct {
	name   string
	ksize  *int64
	srange Range
	weight float64
}

// Returns a pointer to int64 because this value can be unset.
func (b *bucket) GetKSize() *int64 {
	return b.ksize
}

// Generates a new pointer to int64 internally.
func (b *bucket) SetKSize(ksize int64) {
	tmpsize := ksize
	b.ksize = &tmpsize
}

func (b *bucket) GetRange() Range {
	return b.srange
}

func (b *bucket) SetRange(srange Range) {
	b.srange = srange
}

func (b *bucket) GetWeight() float64 {
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

// A StackBucket is exactly like a bucket except that it supports DefaultSize().
type StackBucket interface {
	Bucket
	DefaultSize() MemSize
}

func NewStackBucket(weight float64, srange Range) (StackBucket, error) {
	return newBucket("stack", weight, srange)
}

var (
	jvm_DEFAULT_STACK_SIZE = NewMemSize(mEGA)
)

// The default stacksize (minimum of the range, or the JRE standard default).
func (b *bucket) DefaultSize() MemSize {
	floor := b.srange.Floor()
	if floor.Bytes() == 0 {
		return jvm_DEFAULT_STACK_SIZE
	}
	return floor
}
