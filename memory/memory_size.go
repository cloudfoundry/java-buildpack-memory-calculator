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
	"strconv"
	"strings"
)

// memory_size.go defines the MemSize struct which captures a memory size
// allocation. It understands the normal nK, nM, nG string representations,
// and permits scaling and comparison operations.  The methods are described
// in-line.
// type MemSize is exported, all the methods have receiver *MemSize

type MemSize struct {
	sizeInBytes int64
}

const (
	bYTE = 1
	kILO = 1024 * bYTE
	mEGA = 1024 * kILO
	gIGA = 1024 * mEGA
)

// The empty memory allocation (not nil).
var MS_ZERO *MemSize

func init() {
	MS_ZERO = NewMemSize(0)
}

// Construct a new MemSize object from an int64
func NewMemSize(ms int64) *MemSize {
	return &MemSize{ms}
}

// Construct a new MemSize object from a string description
//
// Errors include:
//	errors from ParseInt
//	error invalid memory size string '%s'
func NewMemSizeFromString(ms string) (*MemSize, error) {
	ms = strings.TrimSpace(ms)
	if ms == "" {
		return nil, fmt.Errorf("memory size string cannot be empty")
	}
	if ms == "0" {
		return NewMemSize(0), nil
	}

	factor, intStr := int64(1), ms[:len(ms)-1]
	switch ms[len(ms)-1] {
	case 'b', 'B':
		factor = bYTE
	case 'k', 'K':
		factor = kILO
	case 'm', 'M':
		factor = mEGA
	case 'g', 'G':
		factor = gIGA
	default:
		return nil, fmt.Errorf("invalid memory size string '%s'", ms)
	}

	num, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return NewMemSize(num * factor), nil
}

// The number of bytes in the MemSize
func (ms *MemSize) Bytes() int64 {
	return ms.sizeInBytes
}

// The number of (whole) kilobytes in the MemSize
func (ms *MemSize) Kilos() int64 {
	return ms.sizeInBytes / kILO
}

// The number of (whole) megabytes in the MemSize
func (ms *MemSize) Megas() int64 {
	return ms.sizeInBytes / mEGA
}

// The number of (whole) gigabytes in the MemSize
func (ms *MemSize) Gigas() int64 {
	return ms.sizeInBytes / gIGA
}

// A string presentation of the MemSize rounded down to whole numbers of
// kilobytes, and expressed in the highest unit that is exact using the K,M,G
// suffices. Less than 1K produces "0" as the string output.
func (ms *MemSize) String() string {
	var (
		val  int64
		suff string
	)

	if v := ms.Kilos(); v == 0 {
		return "0"
	} else if v%mEGA == 0 {
		val, suff = v/mEGA, "G"
	} else if v%kILO == 0 {
		val, suff = v/kILO, "M"
	} else {
		val, suff = v, "K"
	}

	return fmt.Sprintf("%d%s", val, suff)
}

// True if the receiver has less bytes in it than does other.
func (ms *MemSize) LessThan(other *MemSize) bool {
	return ms.Bytes() < other.Bytes()
}

// Produce a new MemSize with the sum of the number of bytes in receiver and other.
func (ms *MemSize) Add(other *MemSize) *MemSize {
	return &MemSize{ms.sizeInBytes + other.sizeInBytes}
}

// Produce a new MemSize with factor times the number of bytes in it (rounded to nearest integer).
func (ms *MemSize) Scale(factor float64) *MemSize {
	return &MemSize{int64(factor*float64(ms.sizeInBytes) + 0.5)}
}

// True if the receiver has exactly the same number of bytes in it as does other.
func (ms *MemSize) Equals(other *MemSize) bool {
	return ms.sizeInBytes == other.sizeInBytes
}

// True if the receiver has exactly zero bytes in it.
func (ms *MemSize) Empty() bool {
	return ms.sizeInBytes == 0
}

// The ratio of the sizes in receiver and other as a floating point number.
// twoGig.DividedBy(oneGig) should return 2.0.
// oneGig.DividedBy(twoGig) should return 0.5.
func (ms *MemSize) DividedBy(other *MemSize) float64 {
	return float64(ms.sizeInBytes) / float64(other.sizeInBytes)
}
