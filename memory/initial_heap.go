/*
 * Copyright 2026 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package memory

import (
	"fmt"
	"regexp"
	"strings"
)

var initialHeapRE = regexp.MustCompile(fmt.Sprintf("^-Xms(%s)$", sizePattern))

type InitialHeap Size

func IsInitialHeap(s string) bool {
	return initialHeapRE.MatchString(strings.TrimSpace(s))
}

func ParseInitialHeap(s string) (InitialHeap, error) {
	t := strings.TrimSpace(s)

	if !initialHeapRE.MatchString(t) {
		return InitialHeap(0), fmt.Errorf("initial heap does not match pattern '%s': %s", initialHeapRE.String(), t)
	}

	groups := initialHeapRE.FindStringSubmatch(t)
	size, err := ParseSize(groups[1])
	if err != nil {
		return InitialHeap(0), err
	}

	return InitialHeap(size), nil
}

func (i InitialHeap) String() string {
	return fmt.Sprintf("-Xms%s", Size(i))
}
