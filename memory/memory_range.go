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

// Range denotes a range of memory sizes. There are both bounded and unbounded
// ranges. A bounded range has both a lower and an upper bound; an unbounded
// range has only a lower bound. Bounds are specified as *MemSize values.
type Range interface {
	Floor() *MemSize                 // The lower bound of the range
	Ceiling() (*MemSize, error)      // The upper bound of the range (returns an error if the range is unbounded)
	IsBounded() bool                 // true iff the range is bounded
	Contains(val *MemSize) bool      // true iff val is in the range
	Constrain(val *MemSize) *MemSize // val if val is in the range, otherwise the nearest bound of the range
	Degenerate() bool                // true iff the lower bound equals the upper bound (always false for unbounded ranges)
	Scale(factor float64) Range      // A new Range with both the lower and upper bounds modified by the given factor. nil if factor is negative.
	Equals(rnge Range) bool          // true iff the range denotes the same range as rnge.

	fmt.Stringer
}

type memRange struct {
	lower     *MemSize
	upper     *MemSize
	unbounded bool
}

// NewRangeFromString produces a Range from a string representation.
//
// strRange ::= numLimit | limit '..' limit
// limit    ::= '' | numLimit
// numLimit ::= '0' | INTEGER unit
// unit     ::= 'b' | 'B' | 'k' | 'K' | 'm' | 'M' | 'g' | 'G'
//
// b|B, k|K, m|M, and g|G denote bytes, kilobytes, megabytes and gigabytes,
// respectively. INTEGER is a positive, negative or zero decimal integer
// string. Space around strRange or either limit is ignored.
//
// The default lower (first) limit is 0, the default upper (second) limit is
// infinity (or, in other words, denotes an unbounded range).
//
// Errors include syntax errors, and 'invalid bounded range' which occurs if
// the lower limit is numerically greater than the upper limit.
func NewRangeFromString(strRange string) (Range, error) {

	bounds := strings.Split(strings.TrimSpace(strRange), "..")

	if len(bounds) == 2 {

		if strings.TrimSpace(bounds[0]) == "" {
			bounds[0] = "0"
		}
		lb, lberr := NewMemSizeFromString(bounds[0])
		if lberr != nil {
			return nil, fmt.Errorf("invalid string range '%s'", strRange)
		}

		if strings.TrimSpace(bounds[1]) == "" {
			return NewUnboundedRange(lb)
		}

		ub, uberr := NewMemSizeFromString(bounds[1])
		if uberr != nil {
			return nil, fmt.Errorf("invalid string range '%s'", strRange)
		}
		return NewRange(lb, ub)

	} else if len(bounds) == 1 {

		lb, lberr := NewMemSizeFromString(bounds[0])
		if lberr != nil {
			return nil, fmt.Errorf("invalid string range '%s'", strRange)
		}
		return NewRange(lb, lb)

	}
	return nil, fmt.Errorf("invalid string range '%s'", strRange)
}

// NewRange produces a bounded range from a lower bound *MemSize and an upper
// bound *MemSize. It is an error if the lower bound is greater than the upper
// bound.
func NewRange(low *MemSize, upp *MemSize) (Range, error) {
	return newRange(low.Bytes(), upp.Bytes(), false)
}

// NewUnboundedRange produces an unbounded range given a lower bound *MemSize.
func NewUnboundedRange(low *MemSize) (Range, error) {
	return newRange(low.Bytes(), 0, true)
}

// Floor() produces the lower bound *MemSize (not a copy)
func (r *memRange) Floor() *MemSize {
	return r.lower
}

// IsBounded() returns true iff the range is bounded.
func (r *memRange) IsBounded() bool {
	return !r.unbounded
}

// Ceiling() produces the upper bound *MemSize (not a copy) if the range is
// bounded. It returns an error if the range is unbounded.
func (r *memRange) Ceiling() (*MemSize, error) {
	if r.unbounded {
		return nil, fmt.Errorf("Cannot take Ceiling() of unbounded range %v...", r.lower)
	}
	return r.upper, nil
}

// Contains(val) returns true iff val is within the range.
func (r *memRange) Contains(val *MemSize) bool {
	return (!val.LessThan(r.lower)) && (r.unbounded || !(r.upper.LessThan(val)))
}

// Constrain(val) returns val (not a copy) iff val is within the range,
// otherwise returns Floor() or Ceiling() whichever is nearer val numerically.
func (r *memRange) Constrain(val *MemSize) *MemSize {
	if val.LessThan(r.lower) {
		return r.lower
	}
	if r.IsBounded() && r.upper.LessThan(val) {
		return r.upper
	}
	return val
}

// Degenerate() returns true iff the range is bounded and the lower bound
// equals the upper bound.
func (r *memRange) Degenerate() bool {
	if r.unbounded {
		return false
	}
	return r.upper.Equals(r.lower)
}

// Scale returns a (new) range which is larger than the range by the factor.
// Unbounded ranges remain unbounded (above). Negative factors result in a nil
// result (which is an invalid Range).
func (r *memRange) Scale(factor float64) Range {
	if factor < 0.0 {
		return nil
	} else {
		nr, _ := newRange(r.lower.Scale(factor).Bytes(), r.upper.Scale(factor).Bytes(), r.unbounded)
		return nr
	}
}

// Returns true iff the range and r2 denote exactly the same range of memory
// sizes.
func (r *memRange) Equals(r2 Range) bool {
	if r.lower.Equals(r2.Floor()) {
		if r.unbounded {
			return !r2.IsBounded()
		}
		if r2.IsBounded() {
			return r.upper.Equals(fstms(r2.Ceiling()))
		}
	}
	return false
}

// Produces a string representation of the range in the same format as
// interpreted by NewRangeFromString(). The limits are String representations
// of the Floor() and Ceiling() *MemSize values, which means that values are
// accurate to the next lower whole kilobyte. limits numerically less than
// 1024 bytes are denoted by '0'.
func (r *memRange) String() string {
	if r.unbounded {
		return fmt.Sprintf("%v..", r.lower)
	}
	return fmt.Sprintf("%v..%v", r.lower, r.upper)
}

func newRange(low, upp int64, unb bool) (Range, error) {
	if !unb && low > upp {
		return nil, fmt.Errorf("invalid bounded range: lower (%d) is higher than upper (%d)", low, upp)
	}

	lowMs := NewMemSize(low)

	var uppMs *MemSize = nil
	if !unb {
		uppMs = NewMemSize(upp)
	}

	return &memRange{
		lower:     lowMs,
		upper:     uppMs,
		unbounded: unb,
	}, nil
}

func fstms(f *MemSize, _ error) *MemSize {
	return f
}
