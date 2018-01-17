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

var _ = Describe("VmOptions", func() {
	const (
		testMemSizeString      = "30M"
		testMemOtherSizeString = "50M"
	)

	var (
		rawOpts          string
		vmOptions        memory.VmOptions
		err              error
		testMemSize      memory.MemSize
		testMemOtherSize memory.MemSize
	)

	BeforeEach(func() {
		testMemSize, err = memory.NewMemSizeFromString(testMemSizeString)
		Ω(err).ShouldNot(HaveOccurred())

		testMemOtherSize, err = memory.NewMemSizeFromString(testMemOtherSizeString)
		Ω(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		vmOptions, err = memory.NewVmOptions(rawOpts)
	})

	Context("when the raw options are empty", func() {
		BeforeEach(func() {
			rawOpts = ""
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when the raw options contain no memory options", func() {
		BeforeEach(func() {
			rawOpts = "-verbose:class"
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should omit the raw options in the delta output", func() {
			Ω(vmOptions.DeltaString()).Should(BeEmpty())
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.MaxPermSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.MaxPermSize)).Should(Equal(testMemSize))
			Ω(vmOptions.DeltaString()).Should(Equal("-XX:MaxPermSize=30M"))
		})
	})

	Context("when the raw options contain maximum heap size", func() {
		BeforeEach(func() {
			rawOpts = "-Xmx" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should omit the raw options from the delta output", func() {
			Ω(vmOptions.DeltaString()).Should(BeEmpty())
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxHeapSize)).Should(Equal(testMemSize))
		})
	})

	Context("when the raw options do not contain maximum heap size", func() {
		BeforeEach(func() {
			rawOpts = ""
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should initially not reproduce the maxmimum heap size in the delta output", func() {
			Ω(vmOptions.DeltaString()).ShouldNot(ContainSubstring("-Xmx"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxHeapSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.MaxHeapSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.MaxHeapSize)).Should(Equal(testMemSize))
			Ω(vmOptions.DeltaString()).Should(ContainSubstring("-Xmx30M"))
		})
	})

	Context("when the raw options contain maximum metaspace size", func() {
		BeforeEach(func() {
			rawOpts = "-XX:MaxMetaspaceSize=" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should omit the raw options from the delta output", func() {
			Ω(vmOptions.DeltaString()).Should(BeEmpty())
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxMetaspaceSize)).Should(Equal(testMemSize))
		})
	})

	Context("when the raw options do not contain maximum metaspace size", func() {
		BeforeEach(func() {
			rawOpts = ""
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should initially not reproduce the maximum metaspace size in the delta output", func() {
			Ω(vmOptions.DeltaString()).ShouldNot(ContainSubstring("MaxMetaspaceSize"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxMetaspaceSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.MaxMetaspaceSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.MaxMetaspaceSize)).Should(Equal(testMemSize))
			Ω(vmOptions.DeltaString()).Should(ContainSubstring("-XX:MaxMetaspaceSize=30M"))
		})
	})

	Context("when the raw options contain maximum permgen size", func() {
		BeforeEach(func() {
			rawOpts = "-XX:MaxPermSize=" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should omit the raw options from the delta output", func() {
			Ω(vmOptions.DeltaString()).Should(BeEmpty())
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxPermSize)).Should(Equal(testMemSize))
		})
	})

	Context("when the raw options do not contain maximum permgen size", func() {
		BeforeEach(func() {
			rawOpts = ""
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should initially not reproduce the maximum permgen size in the delta output", func() {
			Ω(vmOptions.DeltaString()).ShouldNot(ContainSubstring("MaxPermSize"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxPermSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.MaxPermSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.MaxPermSize)).Should(Equal(testMemSize))
			Ω(vmOptions.DeltaString()).Should(ContainSubstring("-XX:MaxPermSize=30M"))
		})
	})

	Context("when the raw options contain stack size", func() {
		BeforeEach(func() {
			rawOpts = "-Xss" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should omit the raw options from the delta output", func() {
			Ω(vmOptions.DeltaString()).Should(BeEmpty())
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.StackSize)).Should(Equal(testMemSize))
		})
	})

	Context("when the raw options do not contain stack size", func() {
		BeforeEach(func() {
			rawOpts = ""
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should initially not reproduce the stack size in the delta output", func() {
			Ω(vmOptions.DeltaString()).ShouldNot(ContainSubstring("Xss"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.StackSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.StackSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.StackSize)).Should(Equal(testMemSize))
			Ω(vmOptions.DeltaString()).Should(ContainSubstring("-Xss30M"))
		})
	})

	Context("when the raw options contain maximum direct memory size", func() {
		BeforeEach(func() {
			rawOpts = "-XX:MaxDirectMemorySize=" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should omit the maximum direct memory size from the delta output", func() {
			Ω(vmOptions.DeltaString()).Should(BeEmpty())
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxDirectMemorySize)).Should(Equal(testMemSize))
		})
	})

	Context("when the raw options do not contain maximum direct memory size", func() {
		BeforeEach(func() {
			rawOpts = ""
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should initially not reproduce the maximum direct memory size in the delta output", func() {
			Ω(vmOptions.DeltaString()).ShouldNot(ContainSubstring("MaxDirectMemorySize"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxDirectMemorySize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.MaxDirectMemorySize, testMemSize)
			Ω(vmOptions.MemOpt(memory.MaxDirectMemorySize)).Should(Equal(testMemSize))
			Ω(vmOptions.DeltaString()).Should(ContainSubstring("-XX:MaxDirectMemorySize=30M"))
		})
	})

	Context("when the raw options contain reserved code cache size", func() {
		BeforeEach(func() {
			rawOpts = "-XX:ReservedCodeCacheSize=" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should omit the reserved code cache size from the delta output", func() {
			Ω(vmOptions.DeltaString()).Should(BeEmpty())
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.ReservedCodeCacheSize)).Should(Equal(testMemSize))
		})
	})

	Context("when the raw options do not contain reserved code cache size", func() {
		BeforeEach(func() {
			rawOpts = ""
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should initially not reproduce the reserved code cache size in the delta output", func() {
			Ω(vmOptions.DeltaString()).ShouldNot(ContainSubstring("ReservedCodeCacheSize"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.ReservedCodeCacheSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.ReservedCodeCacheSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.ReservedCodeCacheSize)).Should(Equal(testMemSize))
			Ω(vmOptions.DeltaString()).Should(ContainSubstring("-XX:ReservedCodeCacheSize=30M"))
		})
	})

	Context("when the raw options contain a duplicated option", func() {
		Context("when the duplicated options are valid", func() {
			BeforeEach(func() {
				rawOpts = "-Xmx" + testMemSizeString + " -Xmx" + testMemOtherSizeString
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should omit the duplicated options from the delta output", func() {
				Ω(vmOptions.DeltaString()).Should(BeEmpty())
			})

			It("should capture the value in the last option", func() {
				Ω(vmOptions.MemOpt(memory.MaxHeapSize)).Should(Equal(testMemOtherSize))
			})
		})

		Context("when an earlier option is invalid", func() {
			BeforeEach(func() {
				rawOpts = "-XmxBAD" + " -Xmx" + testMemOtherSizeString
			})

			It("should return a suitable error", func() {
				Ω(err).Should(MatchError("invalid memory size string 'BAD'"))
			})
		})

		Context("when a later option is invalid", func() {
			BeforeEach(func() {
				rawOpts = "-Xmx" + testMemSizeString + " -XmxBAD"
			})

			It("should return a suitable error", func() {
				Ω(err).Should(MatchError("invalid memory size string 'BAD'"))
			})
		})
	})

	Describe("Copy", func() {
		BeforeEach(func() {
			rawOpts = "-XX:ReservedCodeCacheSize=" + testMemSizeString
		})

		It("should copy raw and set options", func() {
			metaSpaceSize, err := memory.NewMemSizeFromString("45M")
			Ω(err).ShouldNot(HaveOccurred())
			vmOptions.SetMemOpt(memory.MaxMetaspaceSize, metaSpaceSize)

			vmOptionsCopy := vmOptions.Copy()

			Ω(vmOptionsCopy.DeltaString()).Should(Equal("-XX:MaxMetaspaceSize=45M"))
			Ω(vmOptionsCopy.String()).Should(Equal("-XX:ReservedCodeCacheSize=30M, -XX:MaxMetaspaceSize=45M"))
		})
	})

	Describe("ClearMemOpt", func() {
		BeforeEach(func() {
			rawOpts = "-XX:ReservedCodeCacheSize=" + testMemSizeString
		})

		It("should clear options which have been set", func() {
			metaSpaceSize, err := memory.NewMemSizeFromString("45M")
			Ω(err).ShouldNot(HaveOccurred())
			vmOptions.SetMemOpt(memory.MaxMetaspaceSize, metaSpaceSize)

			Ω(vmOptions.DeltaString()).Should(Equal("-XX:MaxMetaspaceSize=45M"))
			Ω(vmOptions.String()).Should(Equal("-XX:ReservedCodeCacheSize=30M, -XX:MaxMetaspaceSize=45M"))

			vmOptions.ClearMemOpt(memory.MaxMetaspaceSize)

			Ω(vmOptions.DeltaString()).Should(Equal(""))
			Ω(vmOptions.String()).Should(Equal("-XX:ReservedCodeCacheSize=30M"))
		})

		It("should clear raw options", func() {
			metaSpaceSize, err := memory.NewMemSizeFromString("45M")
			Ω(err).ShouldNot(HaveOccurred())
			vmOptions.SetMemOpt(memory.MaxMetaspaceSize, metaSpaceSize)

			Ω(vmOptions.String()).Should(Equal("-XX:ReservedCodeCacheSize=30M, -XX:MaxMetaspaceSize=45M"))

			vmOptions.ClearMemOpt(memory.ReservedCodeCacheSize)

			Ω(vmOptions.String()).Should(Equal("-XX:MaxMetaspaceSize=45M"))
		})
	})

	Describe("String", func() {
		BeforeEach(func() {
			rawOpts = "-XX:ReservedCodeCacheSize=" + testMemSizeString
		})

		It("should record any value set", func() {
			metaSpaceSize, err := memory.NewMemSizeFromString("45M")
			Ω(err).ShouldNot(HaveOccurred())
			vmOptions.SetMemOpt(memory.MaxMetaspaceSize, metaSpaceSize)

			Ω(vmOptions.DeltaString()).Should(Equal("-XX:MaxMetaspaceSize=45M"))
			Ω(vmOptions.String()).Should(Equal("-XX:ReservedCodeCacheSize=30M, -XX:MaxMetaspaceSize=45M"))
		})
	})
})
