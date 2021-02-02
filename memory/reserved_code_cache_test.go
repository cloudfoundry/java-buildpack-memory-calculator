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

package memory_test

import (
	"testing"

	"github.com/instana/java-buildpack-memory-calculator/v4/memory"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func TestReservedCodeCache(t *testing.T) {
	spec.Run(t, "ReservedCodeCache", func(t *testing.T, when spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("formats", func() {
			g.Expect(memory.ReservedCodeCache(memory.Kibi).String()).To(Equal("-XX:ReservedCodeCacheSize=1K"))
		})

		it("matches -XX:ReservedCodeCacheSize", func() {
			g.Expect(memory.IsReservedCodeCache("-XX:ReservedCodeCacheSize=1K")).To(BeTrue())
		})

		it("does not match non -XX:ReservedCodeCacheSize", func() {
			g.Expect(memory.IsReservedCodeCache("-Xss1K")).To(BeFalse())
		})

		it("parses", func() {
			g.Expect(memory.ParseReservedCodeCache("-XX:ReservedCodeCacheSize=1K")).To(Equal(memory.ReservedCodeCache(memory.Kibi)))
		})

	})
}
