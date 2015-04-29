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

var _ = Describe("MemoryBucket", func() {
	var (
		DEFAULT_JRE_STACK_SIZE = memory.NewMemSize(mEGA)
		testRange              memory.Range
		testZeroRange          memory.Range
		testUBRange            memory.Range
		shouldFail             func(memory.Bucket, error)
		shouldWork             func(memory.Bucket, error) memory.Bucket
	)

	BeforeEach(func() {
		testRange = getBMR(2*mEGA, 3*mEGA)
		testZeroRange = getBMR(0, 4*mEGA)
		testUBRange = getUMR(10 * kILO)

		shouldFail = func(b memory.Bucket, err error) {
			Ω(b).Should(BeNil())
			Ω(err).Should(HaveOccurred())
		}

		shouldWork = func(b memory.Bucket, err error) memory.Bucket {
			Ω(b).ShouldNot(BeNil())
			Ω(err).ShouldNot(HaveOccurred())
			return b
		}
	})

	Context("constructors", func() {

		Context("work", func() {
			It("with non-blank name and good weights", func() {
				b := shouldWork(memory.NewBucket("abucketname", 0.2, testRange))

				Ω(b.Name()).Should(Equal("abucketname"))
				Ω(b.Range()).Should(Equal(testRange))
				Ω(b.GetSize()).Should(BeNil())
				b.SetSize(124)
				Ω(*b.GetSize()).Should(Equal(memory.MemSize(124)))
				Ω(b.Weight()).Should(BeNumerically("~", 0.2))
				Ω(b.DefaultSize()).Should(Equal(memory.MS_ZERO))
			})

			It("with 'stack' bucket and good weights", func() {
				sb := shouldWork(memory.NewBucket("stack", 0.1, testZeroRange))

				Ω(sb.Name()).Should(Equal("stack"))
				Ω(sb.DefaultSize()).Should(Equal(DEFAULT_JRE_STACK_SIZE))

				sb = shouldWork(memory.NewBucket("stack", 0.1, testRange))
				Ω(sb.DefaultSize()).Should(Equal(memory.NewMemSize(2 * mEGA)))
			})

			It("with spaced non-blank name", func() {
				b := shouldWork(memory.NewBucket("  \t abucketname ", 0.0, testRange))
				Ω(b.Name()).Should(Equal("abucketname"))

				b2 := shouldWork(memory.NewBucket("abucketname", 0.0, testRange))
				Ω(b).Should(Equal(b2))
			})

			It("with non-zero weights", func() {
				b := shouldWork(memory.NewBucket("abucketname", 1.0, testRange))
				Ω(b.Weight()).Should(BeNumerically("~", 1.0))
				b = shouldWork(memory.NewBucket("abucketname", 0.2, testRange))
				Ω(b.Weight()).Should(BeNumerically("~", 0.2))
				b = shouldWork(memory.NewBucket("abucketname", 0.9, testRange))
				Ω(b.Weight()).Should(BeNumerically("~", 0.9))
			})
		})

		Context("fail", func() {
			It("with bad names", func() {
				shouldFail(memory.NewBucket("", 0.0, testRange))
				shouldFail(memory.NewBucket("   ", 0.0, testRange))
				shouldFail(memory.NewBucket("  \t", 0.0, testRange))
			})

			It("with bad weights", func() {
				shouldFail(memory.NewBucket("abucket", -0.01, testRange))
				shouldFail(memory.NewBucket("abucket", 10.0, testRange))
				shouldFail(memory.NewBucket("stack", 1.01, testRange))
			})
		})
	})

	Context("operations", func() {
		It("sets Size correctly", func() {
			b := shouldWork(memory.NewBucket("abucket", 0.1, testRange))
			Ω(b.GetSize()).Should(BeNil())

			b.SetSize(memory.MS_ZERO)
			checkSize(b, memory.MS_ZERO)

			b.SetSize(getMs(3 * mEGA))
			checkSize(b, getMs(3*mEGA))
		})

		It("sets Range correctly", func() {
			b := shouldWork(memory.NewBucket("abucket", 0.1, testRange))
			Ω(b.Range()).Should(Equal(getBMR(2*mEGA, 3*mEGA)))

			b.SetRange(getUMR(10 * kILO))
			Ω(b.Range()).Should(Equal(testUBRange))

			b.SetRange(testZeroRange)
			Ω(b.Range()).Should(Equal(testZeroRange))
		})
	})
})

func checkSize(b memory.Bucket, ms memory.MemSize) {
	Ω(b.GetSize()).ShouldNot(BeNil())
	Ω(*b.GetSize()).Should(Equal(ms))
}
