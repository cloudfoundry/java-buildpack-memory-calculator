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

package flags_test

import (
	"testing"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/flags"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func TestJVMOptions(t *testing.T) {
	spec.Run(t, "JVMOptions", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("it is always valid", func() {
			j := flags.JVMOptions{}

			g.Expect(j.Validate()).To(Succeed())
		})

		it("creates string", func() {
			d := memory.MaxDirectMemory(memory.Kibi)
			h := memory.MaxHeap(memory.Kibi)
			m := memory.MaxMetaspace(memory.Kibi)
			r := memory.ReservedCodeCache(memory.Kibi)
			s := memory.Stack(memory.Kibi)

			j := flags.JVMOptions{MaxDirectMemory: &d, MaxHeap: &h, MaxMetaspace: &m, ReservedCodeCache: &r, Stack: &s}

			g.Expect(j.String()).To(Equal("-XX:MaxDirectMemorySize=1K -Xmx1K -XX:MaxMetaspaceSize=1K -XX:ReservedCodeCacheSize=1K -Xss1K"))
		})

		it("parses value", func() {
			d := memory.MaxDirectMemory(memory.Kibi)
			h := memory.MaxHeap(memory.Kibi)
			m := memory.MaxMetaspace(memory.Kibi)
			r := memory.ReservedCodeCache(memory.Kibi)
			s := memory.Stack(memory.Kibi)

			e := flags.JVMOptions{MaxDirectMemory: &d, MaxHeap: &h, MaxMetaspace: &m, ReservedCodeCache: &r, Stack: &s}

			var j flags.JVMOptions

			g.Expect(j.Set("-XX:MaxDirectMemorySize=1K -Xmx1K -XX:MaxMetaspaceSize=1K -XX:ReservedCodeCacheSize=1K -Xss1K")).To(Succeed())
			g.Expect(j).To(Equal(e))
		})

	})
}
