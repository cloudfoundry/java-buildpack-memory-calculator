// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2015 the original author or authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memory_test

import (
	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MemoryRange", func() {

	Context("string constructor", func() {
		var (
			itWorks func(str string, lo memory.MemSize) memory.Range
			itFails func(str string)
		)

		BeforeEach(func() {
			itWorks = func(str string, lo memory.MemSize) memory.Range {
				rnge, err := memory.NewRangeFromString(str)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(rnge).ShouldNot(BeNil())
				Ω(rnge.Floor()).Should(Equal(lo))
				return rnge
			}

			itFails = func(str string) {
				_, err := memory.NewRangeFromString(str)
				Ω(err).Should(HaveOccurred())
			}
		})

		Context("succeeds", func() {
			It("creates a degenerate range correctly", func() {
				rnge := itWorks(" 3m ", getMs(3*mEGA))
				Ω(rnge.IsBounded()).Should(BeTrue())
				Ω(rnge.Degenerate()).Should(BeTrue())
				Ω(rnge.Ceiling()).Should(Equal(rnge.Floor()))
			})

			It("creates a standard range correctly", func() {
				rnge := itWorks(" 3m .. 5m ", getMs(3*mEGA))
				Ω(rnge.IsBounded()).Should(BeTrue())
				Ω(rnge.Degenerate()).Should(BeFalse())
				Ω(rnge.Ceiling()).Should(Equal(getMs(5 * mEGA)))
			})

			It("creates a range with implicit lower bound correctly", func() {
				rnge := itWorks("  .. 5m ", memory.MEMSIZE_ZERO)
				Ω(rnge).Should(Equal(itWorks("0..5m", memory.MEMSIZE_ZERO)))
			})

			It("creates a range with no upper bound correctly", func() {
				rnge := itWorks(" 5m .. ", getMs(5*mEGA))
				Ω(rnge.IsBounded()).Should(BeFalse())
				Ω(rnge.Degenerate()).Should(BeFalse())
				_, err := rnge.Ceiling()
				Ω(err).Should(HaveOccurred())
			})

			It("creates a range with no upper or lower bound correctly", func() {
				Ω(itWorks("..", memory.MEMSIZE_ZERO)).Should(Equal(itWorks("0..", memory.MEMSIZE_ZERO)))
			})

			It("creates a range of 0.. with an empty range string", func() {
				Ω(itWorks("", memory.MEMSIZE_ZERO)).Should(Equal(itWorks("0..", memory.MEMSIZE_ZERO)))
			})
		})

		Context("fails", func() {
			It("fails to create an empty range", func() {
				itFails("2m..1m")
			})
			It("fails with invalid syntax", func() {
				itFails("1m..2m..4m")
				itFails("1m...2m")
				itFails("1gg..2g")
				itFails("1a..")
			})
		})
	})

	Context("memory_size constructors", func() {
		var (
			boundedWorks   func(lo, hi memory.MemSize) memory.Range
			unboundedWorks func(lo memory.MemSize) memory.Range
			boundedFails   func(lo, hi memory.MemSize)
		)

		BeforeEach(func() {
			boundedWorks = func(lo, hi memory.MemSize) memory.Range {
				rnge, err := memory.NewRange(lo, hi)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(rnge).ShouldNot(BeNil())
				Ω(rnge.Floor()).Should(Equal(lo))
				Ω(rnge.Ceiling()).Should(Equal(hi))
				if lo.Equals(hi) {
					Ω(rnge.Degenerate()).Should(BeTrue())
				}
				return rnge
			}

			unboundedWorks = func(lo memory.MemSize) memory.Range {
				rnge, err := memory.NewUnboundedRange(lo)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(rnge).ShouldNot(BeNil())
				Ω(rnge.Floor()).Should(Equal(lo))
				Ω(rnge.Degenerate()).Should(BeFalse())
				_, err = rnge.Ceiling()
				Ω(err).Should(HaveOccurred())
				return rnge
			}

			boundedFails = func(lo, hi memory.MemSize) {
				_, err := memory.NewRange(lo, hi)
				Ω(err).Should(HaveOccurred())
			}
		})

		It("succeed with MemSize size ranges", func() {
			boundedWorks(getMs(1*mEGA), getMs(2*mEGA))
			unboundedWorks(getMs(1 * mEGA))
		})
		It("fail with invalid MemSize size ranges", func() {
			boundedFails(getMs(2*mEGA), getMs(1*mEGA))
		})
	})

	Context("operations", func() {
		It("detects memory sizes in and outside a bounded range correctly", func() {
			bmr := boundedMemoryRange(mEGA, gIGA)
			Ω(bmr.Contains(getMs(mEGA))).Should(BeTrue())
			Ω(bmr.Contains(getMs(mEGA + 1))).Should(BeTrue())
			Ω(bmr.Contains(getMs(mEGA - 1))).Should(BeFalse())

			Ω(bmr.Contains(getMs(gIGA))).Should(BeTrue())
			Ω(bmr.Contains(getMs(gIGA + 1))).Should(BeFalse())
			Ω(bmr.Contains(getMs(gIGA - 1))).Should(BeTrue())

			Ω(bmr.Contains(getMs(2 * mEGA))).Should(BeTrue())
			Ω(bmr.Contains(getMs(2 * gIGA))).Should(BeFalse())
			Ω(bmr.Contains(getMs(2 * kILO))).Should(BeFalse())
			Ω(bmr.Contains(getMs(-2 * mEGA))).Should(BeFalse())
		})

		It("detects memory sizes in and outside an unbounded range correctly", func() {
			umr := unboundedMemoryRange(mEGA)
			Ω(umr.Contains(getMs(mEGA))).Should(BeTrue())
			Ω(umr.Contains(getMs(mEGA + 1))).Should(BeTrue())
			Ω(umr.Contains(getMs(mEGA - 1))).Should(BeFalse())

			Ω(umr.Contains(getMs(2 * gIGA))).Should(BeTrue())
			Ω(umr.Contains(getMs(2 * kILO))).Should(BeFalse())
			Ω(umr.Contains(getMs(-2 * gIGA))).Should(BeFalse())
		})

		It("constrains memory sizes in and outside a bounded range correctly", func() {
			bmr := boundedMemoryRange(mEGA, gIGA)
			Ω(bmr.Constrain(getMs(mEGA))).Should(Equal(getMs(mEGA)))
			Ω(bmr.Constrain(getMs(mEGA + 1))).Should(Equal(getMs(mEGA + 1)))
			Ω(bmr.Constrain(getMs(mEGA - 1))).Should(Equal(getMs(mEGA)))

			Ω(bmr.Constrain(getMs(gIGA))).Should(Equal(getMs(gIGA)))
			Ω(bmr.Constrain(getMs(gIGA + 1))).Should(Equal(getMs(gIGA)))
			Ω(bmr.Constrain(getMs(gIGA - 1))).Should(Equal(getMs(gIGA - 1)))

			Ω(bmr.Constrain(getMs(2 * gIGA))).Should(Equal(getMs(gIGA)))
			Ω(bmr.Constrain(getMs(2 * mEGA))).Should(Equal(getMs(2 * mEGA)))
			Ω(bmr.Constrain(getMs(2 * kILO))).Should(Equal(getMs(mEGA)))
			Ω(bmr.Constrain(getMs(-2 * mEGA))).Should(Equal(getMs(mEGA)))
		})

		It("constrains memory sizes in and outside an unbounded range correctly", func() {
			umr := unboundedMemoryRange(mEGA)
			Ω(umr.Constrain(getMs(mEGA))).Should(Equal(getMs(mEGA)))
			Ω(umr.Constrain(getMs(mEGA + 1))).Should(Equal(getMs(mEGA + 1)))
			Ω(umr.Constrain(getMs(mEGA - 1))).Should(Equal(getMs(mEGA)))

			Ω(umr.Constrain(getMs(2 * gIGA))).Should(Equal(getMs(2 * gIGA)))
			Ω(umr.Constrain(getMs(2 * kILO))).Should(Equal(getMs(mEGA)))
			Ω(umr.Constrain(getMs(-2 * gIGA))).Should(Equal(getMs(mEGA)))
		})

		It("compares ranges for equality", func() {
			bmr1, err := memory.NewRangeFromString("3m..5m")
			Ω(err).ShouldNot(HaveOccurred())
			bmr2 := boundedMemoryRange(3*mEGA, 5*mEGA)
			Ω(bmr1.Equals(bmr2)).Should(BeTrue())
			Ω(bmr2.Equals(bmr1)).Should(BeTrue())
			bmr3 := boundedMemoryRange(3*mEGA, 6*mEGA)
			Ω(bmr1.Equals(bmr3)).Should(BeFalse())
			Ω(bmr3.Equals(bmr1)).Should(BeFalse())
			Ω(bmr3.Equals(bmr3)).Should(BeTrue())

			umr1 := unboundedMemoryRange(3 * mEGA)
			Ω(umr1.Equals(umr1)).Should(BeTrue())
			Ω(bmr1.Equals(umr1)).Should(BeFalse())
			Ω(umr1.Equals(bmr2)).Should(BeFalse())

			umr2, err := memory.NewRangeFromString("3m..")
			Ω(umr1.Equals(umr2)).Should(BeTrue())
			Ω(umr2.Equals(umr1)).Should(BeTrue())
		})

		It("scales by a factor", func() {
			bmr := boundedMemoryRange(3*mEGA, 5*mEGA)
			Ω(bmr.Scale(2.0)).Should(Equal(boundedMemoryRange(6*mEGA, 10*mEGA)))
			Ω(bmr.Scale(0.5)).Should(Equal(boundedMemoryRange(3*(mEGA/2), 5*(mEGA/2))))
			Ω(bmr.Scale(-1.0)).Should(BeNil())
		})
	})

	It("produces correct string representations", func() {
		Ω(boundedMemoryRange(3*mEGA, 5*mEGA).String()).Should(Equal("3M..5M"))
		Ω(boundedMemoryRange(3*kILO, 5*mEGA).String()).Should(Equal("3K..5M"))
		Ω(boundedMemoryRange(-3*mEGA, 5*mEGA).String()).Should(Equal("-3M..5M"))
		Ω(unboundedMemoryRange(3 * mEGA).String()).Should(Equal("3M.."))
		Ω(unboundedMemoryRange(-3 * mEGA).String()).Should(Equal("-3M.."))
		Ω(unboundedMemoryRange(0).String()).Should(Equal("0.."))
		Ω(boundedMemoryRange(0, 2*gIGA).String()).Should(Equal("0..2G"))
	})
})
