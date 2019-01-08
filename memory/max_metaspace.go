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

var maxMetaspaceRE = regexp.MustCompile(fmt.Sprintf("^-XX:MaxMetaspaceSize=(%s)$", sizePattern))

type MaxMetaspace Size

func IsMaxMetaspace(s string) bool {
	return maxMetaspaceRE.MatchString(strings.TrimSpace(s))
}

func ParseMaxMetaspace(s string) (MaxMetaspace, error) {
	t := strings.TrimSpace(s)

	if !maxMetaspaceRE.MatchString(t) {
		return MaxMetaspace(0), fmt.Errorf("max metaspace does not match pattern '%s': %s", maxMetaspaceRE.String(), t)
	}

	groups := maxMetaspaceRE.FindStringSubmatch(t)
	size, err := ParseSize(groups[1])
	if err != nil {
		return MaxMetaspace(0), err
	}

	return MaxMetaspace(size), nil
}

func (m MaxMetaspace) String() string {
	return fmt.Sprintf("-XX:MaxMetaspaceSize=%s", Size(m))
}
