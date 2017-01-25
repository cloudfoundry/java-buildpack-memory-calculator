// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2015-2017 the original author or authors.
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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

// Path to calculator executable
var jbmcExec string

func TestIntegration(t *testing.T) {

	// build main program for testing
	SynchronizedBeforeSuite(func() []byte {
		jbmcExec, err := gexec.Build("github.com/cloudfoundry/java-buildpack-memory-calculator", "-a", "-race")
		Ω(err).ShouldNot(HaveOccurred())
		return []byte(jbmcExec)
	}, func(jbmcBinPath []byte) {
		jbmcExec = string(jbmcBinPath)
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
