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

var _ = Describe("MemoryBucket and StackMemoryBucket", func() {
	Context("constructors", func() {
		var (
			DEFAULT_JRE_STACK_SIZE = memory.NewMemSize(mEGA)
			testRange              memory.Range
			testZeroRange          memory.Range
			shouldFail             func(memory.Bucket, error)
			shouldWork             func(memory.Bucket, error) memory.Bucket
		)

		BeforeEach(func() {
			var err error
			testRange, err = memory.NewRangeFromString("2m..3m")
			Ω(err).ShouldNot(HaveOccurred())
			testZeroRange, err = memory.NewRangeFromString("0..3m")
			Ω(err).ShouldNot(HaveOccurred())

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

		Context("works", func() {
			It("with non-blank name and good weights", func() {
				b := shouldWork(memory.NewBucket("abucketname", 0.0, testRange))

				Ω(b.Name()).Should(Equal("abucketname"))
				Ω(b.GetRange()).Should(Equal(testRange))
				Ω(b.GetKSize()).Should(BeNil())
				b.SetKSize(124)
				Ω(*b.GetKSize()).Should(Equal(int64(124)))
				Ω(b.GetWeight()).Should(BeNumerically("~", 0.0))
			})

			It("with StackBucket and good weights", func() {
				sb, err := memory.NewStackBucket(0.0, testZeroRange)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(sb.Name()).Should(Equal("stack"))
				Ω(sb.GetRange()).Should(Equal(testZeroRange))
				Ω(sb.GetKSize()).Should(BeNil())
				sb.SetKSize(124)
				Ω(*sb.GetKSize()).Should(Equal(int64(124)))
				Ω(sb.GetWeight()).Should(BeNumerically("~", 0.0))
				Ω(sb.DefaultSize()).Should(Equal(DEFAULT_JRE_STACK_SIZE))

			})

			It("the same with spaced non-blank name", func() {
				b := shouldWork(memory.NewBucket("  \t abucketname ", 0.0, testRange))
				Ω(b.Name()).Should(Equal("abucketname"))

				Ω(shouldWork(memory.NewBucket("abucketname", 0.0, testRange))).Should(Equal(b))
			})

			It("with non-zero weights", func() {
				b := shouldWork(memory.NewBucket("abucketname", 1.0, testRange))
				Ω(b.GetWeight()).Should(BeNumerically("~", 1.0))
				b = shouldWork(memory.NewBucket("abucketname", 0.2, testRange))
				Ω(b.GetWeight()).Should(BeNumerically("~", 0.2))
				b = shouldWork(memory.NewBucket("abucketname", 0.9, testRange))
				Ω(b.GetWeight()).Should(BeNumerically("~", 0.9))
			})
		})

		Context("fails", func() {
			It("fails with blank names", func() {
				shouldFail(memory.NewBucket("", 0.0, testRange))
				shouldFail(memory.NewBucket("   ", 0.0, testRange))
				shouldFail(memory.NewBucket("  \t", 0.0, testRange))
			})

			It("fails with bad weights", func() {
				shouldFail(memory.NewBucket("abucket", -0.01, testRange))
				shouldFail(memory.NewBucket("abucket", 10.0, testRange))
				shouldFail(memory.NewBucket("abucket", 1.01, testRange))
				shouldFail(memory.NewStackBucket(-0.01, testRange))
				shouldFail(memory.NewStackBucket(10.0, testRange))
				shouldFail(memory.NewStackBucket(1.01, testRange))
			})
		})
	})
})
