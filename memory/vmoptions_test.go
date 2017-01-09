// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2016 the original author or authors.
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
		testMemSizeString = "30M"
	)

	var (
		rawOpts     string
		vmOptions   memory.VmOptions
		err         error
		testMemSize memory.MemSize
	)

	BeforeEach(func() {
		testMemSize, err = memory.NewMemSizeFromString(testMemSizeString)
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

		It("should reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).Should(Equal("-verbose:class"))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.CompressedClassSpaceSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.CompressedClassSpaceSize)).Should(Equal(testMemSize))
			Ω(vmOptions.String()).Should(Equal("-verbose:class -XX:CompressedClassSpaceSize=30M"))
		})
	})

	Context("when the raw options contain maximum heap size", func() {
		BeforeEach(func() {
			rawOpts = "-Xmx" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).Should(Equal("-Xmx30M"))
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

		It("should not reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).ShouldNot(ContainSubstring("-Xmx"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxHeapSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.MaxHeapSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.MaxHeapSize)).Should(Equal(testMemSize))
			Ω(vmOptions.String()).Should(ContainSubstring("-Xmx30M"))
		})
	})

	Context("when the raw options contain maximum metaspace size", func() {
		BeforeEach(func() {
			rawOpts = "-XX:MaxMetaspaceSize=" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).Should(Equal("-XX:MaxMetaspaceSize=30M"))
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

		It("should not reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).ShouldNot(ContainSubstring("MaxMetaspaceSize"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxMetaspaceSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.MaxMetaspaceSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.MaxMetaspaceSize)).Should(Equal(testMemSize))
			Ω(vmOptions.String()).Should(ContainSubstring("-XX:MaxMetaspaceSize=30M"))
		})
	})

	Context("when the raw options contain stack size", func() {
		BeforeEach(func() {
			rawOpts = "-Xss" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).Should(Equal("-Xss30M"))
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

		It("should not reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).ShouldNot(ContainSubstring("Xss"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.StackSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.StackSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.StackSize)).Should(Equal(testMemSize))
			Ω(vmOptions.String()).Should(ContainSubstring("-Xss30M"))
		})
	})

	Context("when the raw options contain maximum direct memory size", func() {
		BeforeEach(func() {
			rawOpts = "-XX:MaxDirectMemorySize=" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).Should(Equal("-XX:MaxDirectMemorySize=30M"))
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

		It("should not reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).ShouldNot(ContainSubstring("MaxDirectMemorySize"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.MaxDirectMemorySize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.MaxDirectMemorySize, testMemSize)
			Ω(vmOptions.MemOpt(memory.MaxDirectMemorySize)).Should(Equal(testMemSize))
			Ω(vmOptions.String()).Should(ContainSubstring("-XX:MaxDirectMemorySize=30M"))
		})
	})

	Context("when the raw options contain reserved code cache size", func() {
		BeforeEach(func() {
			rawOpts = "-XX:ReservedCodeCacheSize=" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).Should(Equal("-XX:ReservedCodeCacheSize=30M"))
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

		It("should not reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).ShouldNot(ContainSubstring("ReservedCodeCacheSize"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.ReservedCodeCacheSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.ReservedCodeCacheSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.ReservedCodeCacheSize)).Should(Equal(testMemSize))
			Ω(vmOptions.String()).Should(ContainSubstring("-XX:ReservedCodeCacheSize=30M"))
		})
	})

	Context("when the raw options contain compressed class space size", func() {
		BeforeEach(func() {
			rawOpts = "-XX:CompressedClassSpaceSize=" + testMemSizeString
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).Should(Equal("-XX:CompressedClassSpaceSize=30M"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.CompressedClassSpaceSize)).Should(Equal(testMemSize))
		})
	})

	Context("when the raw options do not contain compressed class space size", func() {
		BeforeEach(func() {
			rawOpts = ""
		})

		It("should not return an error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should not reproduce the raw options in the output", func() {
			Ω(vmOptions.String()).ShouldNot(ContainSubstring("CompressedClassSpaceSize"))
		})

		It("should capture the value in the correct option", func() {
			Ω(vmOptions.MemOpt(memory.CompressedClassSpaceSize)).Should(Equal(memory.MEMSIZE_ZERO))
		})

		It("should record any value set", func() {
			vmOptions.SetMemOpt(memory.CompressedClassSpaceSize, testMemSize)
			Ω(vmOptions.MemOpt(memory.CompressedClassSpaceSize)).Should(Equal(testMemSize))
			Ω(vmOptions.String()).Should(ContainSubstring("-XX:CompressedClassSpaceSize=30M"))
		})
	})

})
