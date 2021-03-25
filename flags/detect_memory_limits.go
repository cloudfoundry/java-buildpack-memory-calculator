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
	"strconv"
)

const (
	DefaultDetectMemoryLimits = DetectMemoryLimits(false)
	FlagDetectMemoryLimits    = "detect-memory-limits"
)

type DetectMemoryLimits bool

func (d *DetectMemoryLimits) Set(s string) error {
	f, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*d = DetectMemoryLimits(f)
	return nil
}

func (d *DetectMemoryLimits) String() string {
	return strconv.FormatBool(bool(*d))
}

func (d *DetectMemoryLimits) Type() string {
	return "bool"
}

func (d *DetectMemoryLimits) Validate() error {
	return nil
}
