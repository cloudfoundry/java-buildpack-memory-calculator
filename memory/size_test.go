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

func TestMemorySize(t *testing.T) {
	spec.Run(t, "Size", func(t *testing.T, when spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		when("format", func() {

			it("formats bytes", func() {
				g.Expect(memory.Size(1023).String()).To(Equal("0"))
			})

			it("formats Kibi", func() {
				g.Expect(memory.Size(memory.Kibi + 1023).String()).To(Equal("1K"))
			})

			it("formats Mibi", func() {
				g.Expect(memory.Size(memory.Mibi + 1023).String()).To(Equal("1M"))
			})

			it("formats Gibi", func() {
				g.Expect(memory.Size(memory.Gibi + 1023).String()).To(Equal("1G"))
			})

			it("formats Tibi", func() {
				g.Expect(memory.Size(memory.Tibi + 1023).String()).To(Equal("1T"))
			})

			it("formats larger than Tibi", func() {
				g.Expect(memory.Size((memory.Tibi * 1024) + 1023).String()).To(Equal("1024T"))
			})
		})

		when("parse", func() {

			it("parses bytes", func() {
				g.Expect(memory.ParseSize("1")).To(Equal(memory.Size(1)))
				g.Expect(memory.ParseSize("1b")).To(Equal(memory.Size(1)))
			})

			it("parses Kibi", func() {
				g.Expect(memory.ParseSize("1k")).To(Equal(memory.Size(memory.Kibi)))
				g.Expect(memory.ParseSize("1K")).To(Equal(memory.Size(memory.Kibi)))
			})

			it("parses Mibi", func() {
				g.Expect(memory.ParseSize("1m")).To(Equal(memory.Size(memory.Mibi)))
				g.Expect(memory.ParseSize("1M")).To(Equal(memory.Size(memory.Mibi)))
			})

			it("parses Gibi", func() {
				g.Expect(memory.ParseSize("1g")).To(Equal(memory.Size(memory.Gibi)))
				g.Expect(memory.ParseSize("1G")).To(Equal(memory.Size(memory.Gibi)))
			})

			it("parses Tibi", func() {
				g.Expect(memory.ParseSize("1t")).To(Equal(memory.Size(memory.Tibi)))
				g.Expect(memory.ParseSize("1T")).To(Equal(memory.Size(memory.Tibi)))
			})

			it("parses zero", func() {
				g.Expect(memory.ParseSize("0")).To(Equal(memory.Size(0)))
			})

			it("trims whitespace", func() {
				g.Expect(memory.ParseSize("\t\r\n 1")).To(Equal(memory.Size(1)))
				g.Expect(memory.ParseSize("1 \t\r\n")).To(Equal(memory.Size(1)))
			})

			it("does not parse empty value", func() {
				_, err := memory.ParseSize("")
				g.Expect(err).To(HaveOccurred())
			})

			it("does not parse negative value", func() {
				_, err := memory.ParseSize("-1")
				g.Expect(err).To(HaveOccurred())
			})

			it("does not parse unknown units", func() {
				_, err := memory.ParseSize("1A")
				g.Expect(err).To(HaveOccurred())
			})

			it("does not parse non-decimal value", func() {
				_, err := memory.ParseSize("0x1")
				g.Expect(err).To(HaveOccurred())
			})

			it("does not parse non-integral value", func() {
				_, err := memory.ParseSize("1.0")
				g.Expect(err).To(HaveOccurred())
			})

			it("does not parse embedded whitespace", func() {
				_, err := memory.ParseSize("1 0")
				g.Expect(err).To(HaveOccurred())
			})
		})
	})
}
