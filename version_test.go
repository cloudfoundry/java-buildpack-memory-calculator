/*
 * Copyright 2026 the original author or authors.
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

package main

import (
	"runtime/debug"
	"testing"

	. "github.com/onsi/gomega"
)

func TestResolveVersion(t *testing.T) {
	tests := []struct {
		name    string
		stamped string
		info    *debug.BuildInfo
		want    string
	}{
		{"stamped version wins over the module version", "4.3.0", &debug.BuildInfo{Main: debug.Module{Version: "v4.2.0"}}, "4.3.0"},
		{"module version is used when nothing is stamped", "", &debug.BuildInfo{Main: debug.Module{Version: "v4.3.0"}}, "v4.3.0"},
		{"development build falls back to dev", "", &debug.BuildInfo{Main: debug.Module{Version: "(devel)"}}, "dev"},
		{"empty module version falls back to dev", "", &debug.BuildInfo{}, "dev"},
		{"missing build info falls back to dev", "", nil, "dev"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			g.Expect(resolveVersion(test.stamped, test.info)).To(Equal(test.want))
		})
	}
}
