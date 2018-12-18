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

package flags

import (
	"fmt"
	"strconv"
)

const (
	DefaultThreadCount = ThreadCount(0)
	FlagThreadCount    = "thread-count"
)

type ThreadCount int

func (t *ThreadCount) Set(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*t = ThreadCount(i)
	return nil
}

func (t *ThreadCount) String() string {
	return strconv.FormatInt(int64(*t), 10)
}

func (t *ThreadCount) Type() string {
	return "int"
}

func (t *ThreadCount) Validate() error {
	if *t == 0 {
		return fmt.Errorf("--%s must be specified", FlagThreadCount)
	}

	if *t < 0 {
		return fmt.Errorf("--%s must be positive: %d", FlagThreadCount, *t)
	}

	return nil
}
