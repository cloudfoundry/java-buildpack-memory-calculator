/*
 * Copyright 2015-2018 the original author or authors.
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
	"strconv"
	"strings"
)

const (
	Kibi = 1024
	Mibi = 1024 * Kibi
	Gibi = 1024 * Mibi
	Tibi = 1024 * Gibi
)

const sizePattern = "[\\d]+[bkmgtBKMGT]?"

var sizeRE = regexp.MustCompile("^([\\d]+)([bkmgtBKMGT]?)$")

type Size int64

func ParseSize(s string) (Size, error) {
	t := strings.TrimSpace(s)

	if !sizeRE.MatchString(t) {
		return Size(0), fmt.Errorf("memory size does not match pattern '%s': %s", sizeRE.String(), t)
	}

	groups := sizeRE.FindStringSubmatch(t)
	size, err := strconv.ParseInt(groups[1], 10, 64)
	if err != nil {
		return Size(0), fmt.Errorf("memory size is not an integer: %s", groups[1])
	}

	switch strings.ToLower(groups[2]) {
	case "k":
		size *= Kibi
	case "m":
		size *= Mibi
	case "g":
		size *= Gibi
	case "t":
		size *= Tibi
	}

	return Size(size), nil
}

func (s Size) String() string {
	b := int64(s) / Kibi

	if b == 0 {
		return "0"
	}

	if b%Gibi == 0 {
		return fmt.Sprintf("%dT", b/Gibi)
	}

	if b%Mibi == 0 {
		return fmt.Sprintf("%dG", b/Mibi)
	}

	if b%Kibi == 0 {
		return fmt.Sprintf("%dM", b/Kibi)
	}

	return fmt.Sprintf("%dK", b)
}
