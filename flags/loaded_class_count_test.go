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

	"github.com/cloudfoundry/java-buildpack-memory-calculator/flags"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func TestLoadedClassCount(t *testing.T) {
	spec.Run(t, "LoadedClassCount", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("is invalid less than 0", func() {
			l := flags.LoadedClassCount(-1)

			g.Expect(l.Validate()).NotTo(Succeed())
		})

		it("is invalid at 0", func() {
			l := flags.LoadedClassCount(0)

			g.Expect(l.Validate()).NotTo(Succeed())
		})

		it("is valid more than 0", func() {
			l := flags.LoadedClassCount(1)

			g.Expect(l.Validate()).To(Succeed())
		})

		it("parses value", func() {
			var l flags.LoadedClassCount

			g.Expect(l.Set("1")).To(Succeed())
			g.Expect(l).To(Equal(flags.LoadedClassCount(1)))
		})
	})
}
