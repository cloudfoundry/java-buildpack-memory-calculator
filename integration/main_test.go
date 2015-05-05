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
package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("java-buildpack-memory-calculator executable", func() {

	It("executes with help on no parms", func() {
		co, err := runOutput()
		Ω(err).Should(HaveOccurred(), jbmcExec)
		Ω(co).Should(ContainSubstring("\njava-buildpack-memory-calculator\n"), "announce name")
		Ω(co).Should(ContainSubstring("-help=false"), "flag prompts")
		Ω(co).Should(ContainSubstring("\nUsage of "), "Usage prefix")
	})

})

func runOutput(args ...string) ([]byte, error) {
	cmd := exec.Command(jbmcExec, args...)
	co, err := cmd.CombinedOutput()
	return co, err
}
