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

package flags

import (
	"fmt"
	"strconv"
)

const (
	DefaultHeapYoungGenerationRatio = HeapYoungGenerationRatio(.3)
	FlagHeapYoungGenerationRatio    = "heap-young-generation-ratio"
)

type HeapYoungGenerationRatio float32

func (d *HeapYoungGenerationRatio) Set(s string) error {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}

	*d = HeapYoungGenerationRatio(f)
	return nil
}

func (d *HeapYoungGenerationRatio) String() string {
	return strconv.FormatFloat(float64(*d), 'f', 2, 32)
}

func (d *HeapYoungGenerationRatio) Type() string {
	return "float32"
}

func (d *HeapYoungGenerationRatio) Validate() error {
	if *d <= 0 || *d >= 1 {
		return fmt.Errorf("--%s must be a valid ration between 0 and 1, extremes excluded: %2f", FlagHeapYoungGenerationRatio, *d)
	}

	return nil
}
