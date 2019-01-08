/*
 * Copyright 2015-2019 the original author or authors.
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

var maxHeapRE = regexp.MustCompile(fmt.Sprintf("^-Xmx(%s)$", sizePattern))

type MaxHeap Size

func IsMaxHeap(s string) bool {
	return maxHeapRE.MatchString(strings.TrimSpace(s))
}

func ParseMaxHeap(s string) (MaxHeap, error) {
	t := strings.TrimSpace(s)

	if !maxHeapRE.MatchString(t) {
		return MaxHeap(0), fmt.Errorf("max heap does not match pattern '%s': %s", maxHeapRE.String(), t)
	}

	groups := maxHeapRE.FindStringSubmatch(t)
	size, err := ParseSize(groups[1])
	if err != nil {
		return MaxHeap(0), err
	}

	return MaxHeap(size), nil
}

func (m MaxHeap) String() string {
	return fmt.Sprintf("-Xmx%s", Size(m))
}
