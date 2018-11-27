// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2015-2018 the original author or authors.
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

var _ = Describe("MemorySize", func() {

	Context("basic constructors", func() {

		var (
			testItWorks func(string, int64)
			testItFails func(string)
		)

		BeforeEach(func() {
			testItWorks = func(memStr string, memVal int64) {
				ms, err := memory.NewMemSizeFromString(memStr)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(ms.Bytes()).Should(Equal(memVal))
			}

			testItFails = func(memStr string) {
				ms, err := memory.NewMemSizeFromString(memStr)
				Ω(ms).Should(BeZero())
				Ω(err).Should(HaveOccurred())
			}
		})

		Context("succeeds", func() {
			It("accepts memory sizes in bytes, kilobytes, megabytes, or gigabytes", func() {
				testItWorks("1024b", kILO)
				testItWorks("1024B", kILO)
				testItWorks("2048B", 2*kILO)
				testItWorks("1m", mEGA)
				testItWorks("1M", mEGA)
				testItWorks("2M", 2*mEGA)
				testItWorks("1g", gIGA)
				testItWorks("1G", gIGA)
				testItWorks("2G", 2*gIGA)
			})

			It("accepts zero (0) as a valid memory size", func() {
				testItWorks("0", 0)
			})

			It("accepts a negative value as a valid memory size", func() {
				testItWorks("-3M", -3*mEGA)
			})

			It("accepts whitespace either side of a valid memory size", func() {
				testItWorks("  1024b\r", kILO)
				testItWorks(" \t\r0  ", 0)
				testItWorks("\t1M    ", mEGA)
			})
		})

		Context("fails", func() {
			It("does not accept 'empty' strings", func() {
				testItFails("")
				testItFails("  ")
				testItFails("\t \r ")
			})

			It("does not accept absent numeric value", func() {
				testItFails("M")
			})

			It("does not accept unqualified non-zero memory sizes", func() {
				testItFails("512")
			})

			It("does not accept unknown units", func() {
				testItFails("512A")
			})

			It("does not accept non-decimal integer values", func() {
				testItFails("0x24b")
			})

			It("does not accept non-integral values", func() {
				testItFails("10.24G")
			})

			It("does not accept embedded whitespace", func() {
				testItFails("10 24G")
			})
		})
	})

	Context("operations", func() {
		var zero, noMeg, oneMeg, twoMeg, oneKilo, twoGig memory.MemSize

		BeforeEach(func() {
			zero = memory.MEMSIZE_ZERO
			noMeg, oneMeg, twoMeg, oneKilo, twoGig = getMs(0), getMs(mEGA), getMs(2*mEGA), getMs(kILO), getMs(2*gIGA)
		})

		It("compares values correctly", func() {
			Ω(oneMeg.LessThan(twoMeg)).To(BeTrue())
			Ω(oneMeg.LessThan(oneMeg)).To(BeFalse())
			Ω(twoMeg.LessThan(oneMeg)).To(BeFalse())
			Ω(twoMeg.LessThan(twoGig)).To(BeTrue())
			Ω(oneKilo.LessThan(oneMeg)).To(BeTrue())
			Ω(oneKilo.LessThan(oneKilo)).To(BeFalse())
			Ω(oneMeg.Equals(oneMeg)).To(BeTrue())
			Ω(noMeg.Equals(noMeg)).To(BeTrue())
			Ω(zero).To(Equal(noMeg))
		})

		It("correctly detects empty cases", func() {
			Ω(oneMeg.Empty()).To(BeFalse())
			Ω(noMeg.Empty()).To(BeTrue())
		})

		It("correctly adds memory sizes", func() {
			Ω(oneMeg.Add(oneMeg)).To(Equal(twoMeg))
			Ω(noMeg.Add(twoGig)).To(Equal(twoGig))
		})

		It("correctly scales negative memory sizes", func() {
			Ω(getMs(-3).Scale(0.5)).To(Equal(getMs(-1)))
		})

		It("correctly scales positive memory sizes", func() {
			Ω(oneMeg.Scale(2.0)).To(Equal(twoMeg))
			Ω(twoGig.Scale(1.0 / 1024)).To(Equal(twoMeg))
			Ω(getMs(3).Scale(0.5)).To(Equal(getMs(2)))
			Ω(getMs(0).Scale(1e8)).To(Equal(memory.MEMSIZE_ZERO))
		})

		It("correctly derives proportion between two memory sizes", func() {
			Ω(twoMeg.DividedBy(oneMeg)).Should(BeNumerically("~", 2.0))
			Ω(oneMeg.DividedBy(twoMeg)).Should(BeNumerically("~", 0.5))
		})

		It("converts a memory size to a string correctly", func() {
			Ω(memory.MEMSIZE_ZERO.String()).Should(Equal("0"))
			Ω(oneMeg.String()).Should(Equal("1M"))
			Ω(oneMeg.Scale(0.75).String()).Should(Equal("768K"))
			Ω(twoMeg.Scale(3.0).Scale(1.5).String()).Should(Equal("9M"))
			Ω(twoMeg.Scale(300.0).Scale(25.6).String()).Should(Equal("15G"))
			Ω(twoGig.Scale(1.0001).Scale(1.0 / 1.0001).String()).Should(Equal("2G"))
		})

		It("converts a memory size to a string by rounding DOWN kilobytes", func() {
			Ω(getMs(1*kILO - 1).String()).Should(Equal("0"))
			Ω(getMs(2*kILO - 1).String()).Should(Equal("1K"))
			Ω(getMs(1*mEGA - 1).String()).Should(Equal("1023K"))
			Ω(getMs(2*mEGA - 1).String()).Should(Equal("2047K"))
			Ω(getMs(1*gIGA - 1).String()).Should(Equal("1048575K"))
			Ω(getMs(2*gIGA + 1023).String()).Should(Equal("2G"))
		})
	})
})
