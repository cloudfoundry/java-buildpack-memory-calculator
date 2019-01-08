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

package flags

import (
	"fmt"
	"strconv"
)

const (
	DefaultLoadedClassCount = LoadedClassCount(0)
	FlagLoadedClassCount    = "loaded-class-count"
)

type LoadedClassCount int

func (l *LoadedClassCount) Set(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*l = LoadedClassCount(i)
	return nil
}

func (l *LoadedClassCount) String() string {
	return strconv.FormatInt(int64(*l), 10)
}

func (l *LoadedClassCount) Type() string {
	return "int"
}

func (l *LoadedClassCount) Validate() error {
	if *l == 0 {
		return fmt.Errorf("--%s must be specified", FlagLoadedClassCount)
	}

	if *l < 0 {
		return fmt.Errorf("--%s must be positive: %d", FlagLoadedClassCount, *l)
	}

	return nil
}
