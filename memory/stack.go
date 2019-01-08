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

const DefaultStack = Stack(Mibi)

var stackRE = regexp.MustCompile(fmt.Sprintf("^-Xss(%s)$", sizePattern))

type Stack Size

func IsStack(s string) bool {
	return stackRE.MatchString(strings.TrimSpace(s))
}

func ParseStack(s string) (Stack, error) {
	t := strings.TrimSpace(s)

	if !stackRE.MatchString(t) {
		return Stack(0), fmt.Errorf("stack size does not match pattern '%s': %s", stackRE.String(), t)
	}

	groups := stackRE.FindStringSubmatch(t)
	size, err := ParseSize(groups[1])
	if err != nil {
		return Stack(0), err
	}

	return Stack(size), nil
}

func (s Stack) String() string {
	return fmt.Sprintf("-Xss%s", Size(s))
}
