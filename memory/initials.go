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

	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory/switches"
)

//Set of initial minimums
var InitialMinimums = map[string]MemSize {
  "heap": MemSize(2097152), //2MB minimum works for Java 7 and Java 8
  "metaspace" : NewMemSize(262144), //256KB
  "permgen" : NewMemSize(1048576), // 1MB
}

func InitialsSwitches(initials map[string]float64, sizes map[string]MemSize, sfs switches.Funs) (strs []string, warnings []string) {
	strs = make([]string, 0, 10)
	warnings = make([]string, 0, 10)
	for initialName, initial := range initials {
		size, ok := sizes[initialName]
		if !ok {
			continue
		}
		initialSize := size.Scale(initial)
		min, ok := InitialMinimums[initialName]
		if !ok {
			min = MEMSIZE_ZERO
		}
		if initialSize.LessThan(min) {
			if size.LessThan(min) {
				//This only happens if there is not enough max memory to work with.  Don't make things worse.
				initialSize = size
			} else {
				warnings = append(warnings, fmt.Sprintf(
					"The configured initial memory size %[1]s for %[2]s is less than the jvm minimum %[3]s.  Setting initial value to %[3]s.",
					initialSize, initialName, min))
				initialSize = min
			}
		}
		strs = append(strs, sfs.Apply(initialName, initialSize.String())...)
	}
	return
}