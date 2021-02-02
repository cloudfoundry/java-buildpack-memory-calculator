/*
 * Copyright 2021 the original author or authors.
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

var maxHeapYoungGenRE = regexp.MustCompile(fmt.Sprintf("^-Xmn(%s)$", sizePattern))

type MaxHeapYoungGeneration Size

func IsMaxHeapYoungGeneration(s string) bool {
	return maxHeapYoungGenRE.MatchString(strings.TrimSpace(s))
}

func ParseMaxHeapYoungGeneration(s string) (MaxHeapYoungGeneration, error) {
	t := strings.TrimSpace(s)

	if !maxHeapYoungGenRE.MatchString(t) {
		return MaxHeapYoungGeneration(0), fmt.Errorf("max heap young generation does not match pattern '%s': %s", maxHeapYoungGenRE.String(), t)
	}

	groups := maxHeapYoungGenRE.FindStringSubmatch(t)
	size, err := ParseSize(groups[1])
	if err != nil {
		return MaxHeapYoungGeneration(0), err
	}

	return MaxHeapYoungGeneration(size), nil
}

func (m MaxHeapYoungGeneration) String() string {
	return fmt.Sprintf("-Xmn%s", Size(m))
}
