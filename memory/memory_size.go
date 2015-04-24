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

type MemSize struct {
	sizeInBytes int64
}

const (
	bYTE = 1
	kILO = 1024 * bYTE
	mEGA = 1024 * kILO
	gIGA = 1024 * mEGA
)

var MS_ZERO *MemSize

func init() {
	MS_ZERO, _ = NewMemSize("0")
}

func NewMemSize(ms string) (*MemSize, error) {
	ms = strings.TrimSpace(ms)
	var bytes int64 = 0
	if ms != "0" {
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

		bytes = num * factor
	}
	return &MemSize{bytes}, nil
}

func (ms *MemSize) Bytes() int64 {
	return ms.sizeInBytes
}

func (ms *MemSize) Kilos() int64 {
	return ms.sizeInBytes / kILO
}

func (ms *MemSize) Megas() int64 {
	return ms.sizeInBytes / mEGA
}

func (ms *MemSize) Gigas() int64 {
	return ms.sizeInBytes / gIGA
}

func (ms *MemSize) String() string {
	var (
		val  int64
		suff string
	)
	if v := ms.Gigas(); v > 0 {
		val, suff = v, "G"
	} else if v := ms.Megas(); v > 0 {
		val, suff = v, "M"
	} else if v := ms.Kilos(); v > 0 {
		val, suff = v, "K"
	} else {
		return "0"
	}
	return fmt.Sprintf("%d%s", val, suff)
}

func (ms *MemSize) LessThan(other *MemSize) bool {
	return ms.Bytes() < other.Bytes()
}

func (ms *MemSize) Add(other *MemSize) *MemSize {
	return &MemSize{ms.sizeInBytes + other.sizeInBytes}
}

func (ms *MemSize) Scale(factor float64) *MemSize {
	return &MemSize{int64(factor*float64(ms.sizeInBytes) + 0.5)}
}

func (ms *MemSize) Equals(other *MemSize) bool {
	return ms.sizeInBytes == other.sizeInBytes
}

func (ms *MemSize) Empty() bool {
	return ms.sizeInBytes == 0
}

func (ms *MemSize) DividedBy(other *MemSize) float64 {
	return float64(ms.sizeInBytes) / float64(other.sizeInBytes)
}
