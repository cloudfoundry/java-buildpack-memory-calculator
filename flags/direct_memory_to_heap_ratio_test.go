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

package flags_test

import (
	"testing"

	"github.com/instana/java-buildpack-memory-calculator/v4/flags"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func TestDirectMemoryToHeapRatio(t *testing.T) {
	spec.Run(t, "DirectMemoryToHeapRatio", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("is invalid less than 0", func() {
			h := flags.DirectMemoryToHeapRatio(-1)

			g.Expect(h.Validate()).NotTo(Succeed())
		})

		it("is invalid equal to 0", func() {
			h := flags.DirectMemoryToHeapRatio(0)

			g.Expect(h.Validate()).NotTo(Succeed())
		})

		it("is invalid more than 1", func() {
			h := flags.DirectMemoryToHeapRatio(1.1)

			g.Expect(h.Validate()).NotTo(Succeed())
		})

		it("is invalid equal to 1", func() {
			h := flags.DirectMemoryToHeapRatio(1)

			g.Expect(h.Validate()).NotTo(Succeed())
		})

		it("is valid between 0 and 1", func() {
			h := flags.DirectMemoryToHeapRatio(0.5)

			g.Expect(h.Validate()).To(Succeed())
		})

		it("parses value", func() {
			var h flags.DirectMemoryToHeapRatio

			g.Expect(h.Set("0.5")).To(Succeed())
			g.Expect(h).To(Equal(flags.DirectMemoryToHeapRatio(.5)))
		})

	})
}
