package memory_test

import (
	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Allocator", func() {

	Context("constructor", func() {

		Context("works", func() {
			var (
				sizes      strmap
				weights    floatmap
				shouldWork func(memory.Allocator, error)
			)

			BeforeEach(func() {
				sizes = strmap{}
				weights = floatmap{
					"heap":    5.0,
					"permgen": 3.0,
					"stack":   1.0,
				}

				shouldWork = func(a memory.Allocator, err error) {
					Ω(a).ShouldNot(BeNil())
					Ω(err).ShouldNot(HaveOccurred())
				}
			})

			It("with normal parameters", func() {
				shouldWork(memory.NewAllocator(sizes, weights))

			})
		})

	})
})

type strmap map[string]string

type floatmap map[string]float64
