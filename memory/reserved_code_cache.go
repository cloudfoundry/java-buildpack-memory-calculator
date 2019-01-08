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

const DefaultReservedCodeCache = ReservedCodeCache(240 * Mibi)

var reservedCodeCacheRE = regexp.MustCompile(fmt.Sprintf("^-XX:ReservedCodeCacheSize=(%s)$", sizePattern))

type ReservedCodeCache Size

func IsReservedCodeCache(s string) bool {
	return reservedCodeCacheRE.MatchString(strings.TrimSpace(s))
}

func ParseReservedCodeCache(s string) (ReservedCodeCache, error) {
	t := strings.TrimSpace(s)

	if !reservedCodeCacheRE.MatchString(t) {
		return ReservedCodeCache(0), fmt.Errorf("reserved code cache does not match pattern '%s': %s", reservedCodeCacheRE.String(), t)
	}

	groups := reservedCodeCacheRE.FindStringSubmatch(t)
	size, err := ParseSize(groups[1])
	if err != nil {
		return ReservedCodeCache(0), err
	}

	return ReservedCodeCache(size), nil
}

func (r ReservedCodeCache) String() string {
	return fmt.Sprintf("-XX:ReservedCodeCacheSize=%s", Size(r))
}
