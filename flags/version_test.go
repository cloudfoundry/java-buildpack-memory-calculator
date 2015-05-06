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

package flags_test

import (
	. "github.com/cloudfoundry/java-buildpack-memory-calculator/flags"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {

	Context("when constructing a Version", func() {
		It("accepts correct syntax", func() {
			itWorks("1", 1, 0, 0)
			itWorks("1.0", 1, 0, 0)
			itWorks("1.0.0", 1, 0, 0)
			itWorks("2", 2, 0, 0)
			itWorks("1.2.3", 1, 2, 3)
			itWorks("101.102", 101, 102, 0)
			itWorks("0.0.10", 0, 0, 10)
			itWorks("1.-0.0", 1, 0, 0)
		})
		It("rejects bad syntax", func() {
			itFails("")
			itFails(" 1")
			itFails(" 1.2.3")
			itFails("1 .2.3")
			itFails("1. 2.3")
			itFails("1.2 .3")
			itFails("1.2. 3")
			itFails("1.2.3 ")
			itFails("1.0.2.")
			itFails("1.0.")
			itFails("1.0.0.0")
			itFails("1.-10.0")
			itFails("1.Ø.0")
			itFails("1.O.O")
		})
	})

	Context("when comparing Versions", func() {
		It("compares them correctly", func() {
			Ω(Version{0, 0, 0}.LessThan(Version{0, 0, 0})).Should(BeFalse())
			Ω(Version{1, 0, 0}.LessThan(Version{1, 0, 0})).Should(BeFalse())
			Ω(Version{1, 1, 0}.LessThan(Version{1, 1, 0})).Should(BeFalse())
			Ω(Version{1, 1, 1}.LessThan(Version{1, 1, 1})).Should(BeFalse())

			Ω(Version{0, 0, 0}.LessThan(Version{0, 0, 1})).Should(BeTrue())
			Ω(Version{0, 0, 1}.LessThan(Version{0, 1, 0})).Should(BeTrue())
			Ω(Version{0, 1, 1}.LessThan(Version{1, 0, 0})).Should(BeTrue())

			Ω(Version{1, 1, 2}.LessThan(Version{1, 1, 1})).Should(BeFalse())
			Ω(Version{1, 2, 1}.LessThan(Version{1, 1, 2})).Should(BeFalse())
			Ω(Version{2, 1, 1}.LessThan(Version{1, 2, 2})).Should(BeFalse())
		})
	})
})

func itWorks(vstr string, a, b, c int) Version {
	nv, err := NewVersion(vstr)
	Ω(err).ShouldNot(HaveOccurred(), "Version "+vstr)
	Ω(nv).Should(Equal(Version{a, b, c}))
	return nv
}

func itFails(vstr string) {
	nv, err := NewVersion(vstr)
	Ω(err).Should(HaveOccurred(), "Version "+vstr)
	Ω(nv).Should(Equal(Version{0, 0, 0}))
}
