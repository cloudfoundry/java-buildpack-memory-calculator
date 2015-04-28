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

import "fmt"

type switchFun func(v string) []string

func apply(f ...string) func(string) []string {
	return func(v string) []string {
		var res = make([]string, 2)
		for _, form := range f {
			res = append(res, fmt.Sprintf(form, v))
		}
		return res
	}
}

var switchFuns = map[string]switchFun{
	"heap":      apply("-Xmx%s", "-Xms%s"),
	"metaspace": apply("-XX:MaxMetaspaceSize=%s", "-XX:MetaspaceSize=%s"),
	"permgen":   apply("-XX:MaxPermSize=%s", "-XX:PermSize=%s"),
	"stack":     apply("-Xss%s"),
}

type Allocator interface {
	Balance(MemSize)  // Balance allocations to buckets within MemSize memory limit
	LowerBounds()     // Allocate without a memory limit
	Switches() string // Generate JRE switches from current allocations
}

type allocator struct {
	buckets map[string]Bucket
}

func NewAllocator(sizes, heuristics map[string]string) *allocator {
	return nil
}

func (a *allocator) Balance(ms MemSize) {

}

func (a *allocator) LowerBounds() {

}

func (a *allocator) Switches() string {
	return ""
}
