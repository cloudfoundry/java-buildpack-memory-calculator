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

	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
)

const (
	DefaultTotalMemory = TotalMemory(0)
	FlagTotalMemory    = "total-memory"
)

const min = memory.Size(1024)

type TotalMemory memory.Size

func (t *TotalMemory) Set(s string) error {
	m, err := memory.ParseSize(s)
	if err != nil {
		return err
	}

	*t = TotalMemory(m)
	return nil
}

func (t *TotalMemory) String() string {
	return memory.Size(*t).String()
}

func (t *TotalMemory) Type() string {
	return "int64"
}

func (t *TotalMemory) Validate() error {
	if *t == 0 {
		return fmt.Errorf("--%s must be specified", FlagTotalMemory)
	}

	if *t < 0 {
		return fmt.Errorf("--%s must be positive: %d", FlagTotalMemory, *t)
	}

	if memory.Size(*t) < min {
		return fmt.Errorf("--%s must be greater than %d: %d", FlagTotalMemory, min, *t)
	}

	return nil
}
