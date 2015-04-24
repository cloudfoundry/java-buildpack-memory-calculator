package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	It("addition is acurate", func() {
		Expect(1 + 1).Should(Equal(2))
	})
})
