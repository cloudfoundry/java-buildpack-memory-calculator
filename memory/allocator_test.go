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
	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory/switches"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type strmap map[string]string

type floatmap map[string]float64

var _ = Describe("Allocator", func() {

	var (
		a       memory.Allocator
		sizes   strmap
		weights floatmap

		shouldWork func(memory.Allocator, error) memory.Allocator
		shouldFail func(memory.Allocator, error)
	)

	BeforeEach(func() {
		sizes = strmap{}

		weights = floatmap{}

		shouldWork = func(a memory.Allocator, err error) memory.Allocator {
			Ω(a).ShouldNot(BeNil())
			Ω(err).ShouldNot(HaveOccurred())
			return a
		}

		shouldFail = func(a memory.Allocator, err error) {
			Ω(a).Should(BeNil())
			Ω(err).Should(HaveOccurred())
		}
	})

	JustBeforeEach(func() {
		a = shouldWork(memory.NewAllocator(sizes, weights))
	})

	Context("constructor", func() {

		Context("with good parameters", func() {

			BeforeEach(func() {
				sizes = strmap{
					"stack":   "2m",
					"heap":    "30m..",
					"permgen": "10m",
				}

				weights = floatmap{
					"stack":   1.0,
					"heap":    5.0,
					"permgen": 3.0,
					"native":  1.0,
				}
			})

			It("succeeds", func() {
				Ω(memory.GetBuckets(a)).Should(ConsistOf(
					"Bucket{name: stack, size: <nil>, range: 2M..2M, weight: 1}",
					"Bucket{name: heap, size: <nil>, range: 30M.., weight: 5}",
					"Bucket{name: permgen, size: <nil>, range: 10M..10M, weight: 3}",
					"Bucket{name: native, size: <nil>, range: 0.., weight: 1}",
				))
			})

			It("sets lower bounds and reports switches correctly", func() {
				a.SetLowerBounds()

				Ω(memory.GetBuckets(a)).Should(ConsistOf(
					"Bucket{name: stack, size: 2M, range: 2M..2M, weight: 1}",
					"Bucket{name: heap, size: 30M, range: 30M.., weight: 5}",
					"Bucket{name: permgen, size: 10M, range: 10M..10M, weight: 3}",
					"Bucket{name: native, size: 0, range: 0.., weight: 1}",
				))

				sws := a.Switches(switches.AllJreSwitchFuns)
				Ω(sws).Should(ConsistOf(
					"-Xmx30M",
					"-Xms30M",
					"-XX:MaxPermSize=10M",
					"-XX:PermSize=10M",
					"-Xss2M",
				)) // heap, permgen, stack
			})
		})

		Context("good balancing", func() {
			var (
				memLimit = memory.MemSize(0)
			)

			JustBeforeEach(func() {
				Ω(a.Balance(memLimit)).ShouldNot(HaveOccurred())
			})

			Context("with single bucket to 'balance'", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "0.."}
					weights = floatmap{"heap": 5.0}
					memLimit = memory.NewMemSize(1024 * mEGA)
				})

				It("fills the bucket up", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: heap, size: 1G, range: 0.., weight: 5}",
					))
				})
			})

			Context("with no buckets to 'balance'", func() {

				BeforeEach(func() {
					sizes = strmap{}
					weights = floatmap{}
					memLimit = memory.NewMemSize(1024 * mEGA)
				})

				It("results in no buckets", func() {
					Ω(memory.GetBuckets(a)).Should(BeEmpty())
				})
			})

			Context("with two buckets to balance", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "0..", "hope": "0.."}
					weights = floatmap{"heap": 1.0, "hope": 3.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("fills the buckets proportionally", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: heap, size: 1G, range: 0.., weight: 1}",
						"Bucket{name: hope, size: 3G, range: 0.., weight: 3}",
					))
				})
			})

			Context("with two buckets to balance with tight limit", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "1G", "hope": "0.."}
					weights = floatmap{"heap": 1.0, "hope": 3.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("fills the remaining bucket proportionally", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: heap, size: 1G, range: 1G..1G, weight: 1}",
						"Bucket{name: hope, size: 3G, range: 0.., weight: 3}",
					))
				})
			})

			Context("with two buckets to balance with very loose limits", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "512M..2048M", "hope": "0.."}
					weights = floatmap{"heap": 1.0, "hope": 3.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("fills the buckets proportionally", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: heap, size: 1G, range: 512M..2G, weight: 1}",
						"Bucket{name: hope, size: 3G, range: 0.., weight: 3}",
					))
				})
			})

			Context("with two buckets to balance with restricting upper limit on one", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "0..512M", "hope": "0.."}
					weights = floatmap{"heap": 1.0, "hope": 3.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("fills the buckets skewed", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: heap, size: 512M, range: 0..512M, weight: 1}",
						"Bucket{name: hope, size: 3584M, range: 0.., weight: 3}",
					))
				})
			})

			Context("with two buckets to balance with restricting lower limit on one", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "2G..", "hope": "0.."}
					weights = floatmap{"heap": 1.0, "hope": 3.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("fills the buckets skewed", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: heap, size: 2G, range: 2G.., weight: 1}",
						"Bucket{name: hope, size: 2G, range: 0.., weight: 3}",
					))
				})
			})

			Context("defaults maximum heap size and permgen size according to the configured weightings", func() {

				BeforeEach(func() {
					sizes = strmap{}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(1024 * mEGA)
				})

				It("fills the bucket up", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: stack, size: 1M, range: 0.., weight: 1}",
						"Bucket{name: heap, size: 512M, range: 0.., weight: 5}",
						"Bucket{name: permgen, size: 314572K, range: 0.., weight: 3}",
						"Bucket{name: native, size: 104857K, range: 0.., weight: 1}",
					))
				})
			})

			Context("with a smallish memory limit", func() {

				BeforeEach(func() {
					sizes = strmap{}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(20 * mEGA)
				})

				It("still defaults the stacksize", func() {
					Ω(memory.GetBuckets(a)).Should(ContainElement(
						"Bucket{name: stack, size: 1M, range: 0.., weight: 1}",
					))
				})
			})

			Context("when maximum heap size is specified", func() {

				BeforeEach(func() {
					sizes = strmap{"stack": "1m", "heap": "3g"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("balances the permgen", func() {
					bucks := memory.GetBuckets(a)
					Ω(bucks).Should(ContainElement(
						"Bucket{name: permgen, size: 471859K, range: 0.., weight: 3}",
					))
					Ω(bucks).Should(ContainElement(
						"Bucket{name: heap, size: 3G, range: 3G..3G, weight: 5}",
					))
				})
			})

			Context("when maximum permgen size is specified", func() {

				BeforeEach(func() {
					sizes = strmap{"stack": "1M", "permgen": "2g"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("balances the heap", func() {
					bucks := memory.GetBuckets(a)
					Ω(bucks).Should(ContainElement(
						"Bucket{name: permgen, size: 2G, range: 2G..2G, weight: 3}",
					))
					Ω(bucks).Should(ContainElement(
						"Bucket{name: heap, size: 1398101K, range: 0.., weight: 5}",
					))
				})
			})

			Context("when thread stack size is specified", func() {

				BeforeEach(func() {
					sizes = strmap{"stack": "2M"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("balances the heap and permgen", func() {
					bucks := memory.GetBuckets(a)
					Ω(bucks).Should(ContainElement(
						"Bucket{name: permgen, size: 1258291K, range: 0.., weight: 3}",
					))
					Ω(bucks).Should(ContainElement(
						"Bucket{name: heap, size: 2G, range: 0.., weight: 5}",
					))
				})
			})

			Context("when thread stack size is specified as a range", func() {

				BeforeEach(func() {
					sizes = strmap{"stack": "2M..3m"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("balances the heap and permgen", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: permgen, size: 1258291K, range: 0.., weight: 3}",
						"Bucket{name: heap, size: 2G, range: 0.., weight: 5}",
						"Bucket{name: stack, size: 2M, range: 2M..3M, weight: 1}",
						"Bucket{name: native, size: 419430K, range: 0.., weight: 1}",
					))
				})
			})

			Context("when thread stack size is specified as a range which impinges on heap and permgen", func() {

				BeforeEach(func() {
					sizes = strmap{"stack": "1g..2g"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("balances the heap and permgen", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: permgen, size: 1G, range: 0.., weight: 3}",
						"Bucket{name: heap, size: 1747626K, range: 0.., weight: 5}",
						"Bucket{name: stack, size: 1G, range: 1G..2G, weight: 1}",
						"Bucket{name: native, size: 349525K, range: 0.., weight: 1}",
					))
				})
			})

			Context("when heap size and permgen size allow for excess memory", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "50m", "permgen": "50m", "stack": "400m..500m"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("balances the heap and permgen", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: permgen, size: 50M, range: 50M..50M, weight: 3}",
						"Bucket{name: heap, size: 50M, range: 50M..50M, weight: 5}",
						"Bucket{name: stack, size: 500M, range: 400M..500M, weight: 1}",
						"Bucket{name: native, size: 3484M, range: 0.., weight: 1}",
					))
				})
			})

			Context("when heap size and permgen size allow for just enough excess memory", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "3000m", "permgen": "196m", "stack": "400m..500m"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("balances the heap and permgen", func() {
					Ω(memory.GetBuckets(a)).Should(ContainElement(
						"Bucket{name: stack, size: 450000K, range: 400M..500M, weight: 1}",
					))
				})
			})

			Context("when heap size and permgen size allow for just enough excess memory", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "1m", "permgen": "1m", "stack": "2m"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("balances the heap and permgen", func() {
					Ω(memory.GetBuckets(a)).Should(ConsistOf(
						"Bucket{name: permgen, size: 1M, range: 1M..1M, weight: 3}",
						"Bucket{name: heap, size: 1M, range: 1M..1M, weight: 5}",
						"Bucket{name: stack, size: 2M, range: 2M..2M, weight: 1}",
						"Bucket{name: native, size: 3772825K, range: 0.., weight: 1}",
					))
				})
			})

			Context("when the specified maximum memory sizes imply the total memory size may be too large", func() {

				BeforeEach(func() {
					sizes = strmap{"heap": "800m", "permgen": "800m"}
					weights = floatmap{"heap": 5.0, "permgen": 3.0, "stack": 1.0, "native": 1.0}
					memLimit = memory.NewMemSize(4 * gIGA)
				})

				It("sets the heap and permgen", func() {
					bucks := memory.GetBuckets(a)
					Ω(bucks).Should(ContainElement(
						"Bucket{name: permgen, size: 800M, range: 800M..800M, weight: 3}"))
					Ω(bucks).Should(ContainElement(
						"Bucket{name: heap, size: 800M, range: 800M..800M, weight: 5}"))
					Ω(a.GetWarnings()).Should(ConsistOf("There is more than 3 times more spare native memory than the default so configured Java memory may be too small or available memory may be too large"))
				})

			})
		})
	})
})
