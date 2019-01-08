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

const DefaultMaxDirectMemory = MaxDirectMemory(10 * Mibi)

var maxDirectMemoryRE = regexp.MustCompile(fmt.Sprintf("^-XX:MaxDirectMemorySize=(%s)$", sizePattern))

type MaxDirectMemory Size

func IsMaxDirectMemory(s string) bool {
	return maxDirectMemoryRE.MatchString(strings.TrimSpace(s))
}

func ParseMaxDirectMemory(s string) (MaxDirectMemory, error) {
	t := strings.TrimSpace(s)

	if !maxDirectMemoryRE.MatchString(t) {
		return MaxDirectMemory(0), fmt.Errorf("max direct memory does not match pattern '%s': %s", maxDirectMemoryRE.String(), t)
	}

	groups := maxDirectMemoryRE.FindStringSubmatch(t)
	size, err := ParseSize(groups[1])
	if err != nil {
		return MaxDirectMemory(0), err
	}

	return MaxDirectMemory(size), nil
}

func (m MaxDirectMemory) String() string {
	return fmt.Sprintf("-XX:MaxDirectMemorySize=%s", Size(m))
}
