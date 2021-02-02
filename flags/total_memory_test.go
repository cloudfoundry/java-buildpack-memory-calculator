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

func TestTotalMemory(t *testing.T) {
	spec.Run(t, "TotalMemory", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("is invalid less than 0", func() {
			t := flags.TotalMemory(-1)

			g.Expect(t.Validate()).NotTo(Succeed())
		})

		it("is invalid at 0", func() {
			t := flags.TotalMemory(0)

			g.Expect(t.Validate()).NotTo(Succeed())
		})

		it("is invalid less than 1K", func() {
			t := flags.TotalMemory(1023)

			g.Expect(t.Validate()).NotTo(Succeed())
		})

		it("is valid more than 1K", func() {
			t := flags.TotalMemory(1024)

			g.Expect(t.Validate()).To(Succeed())
		})

		it("parses value", func() {
			var t flags.TotalMemory

			g.Expect(t.Set("1K")).To(Succeed())
			g.Expect(t).To(Equal(flags.TotalMemory(1024)))
		})
	})
}
