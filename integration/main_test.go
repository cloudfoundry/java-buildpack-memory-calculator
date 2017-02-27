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
	"bytes"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("java-buildpack-memory-calculator executable", func() {

	It("executes with help and usage on no parms", func() {
		co, err := runOutput()
		Ω(err).Should(HaveOccurred(), jbmcExec)
		Ω(co).Should(ContainSubstring("\njava-buildpack-memory-calculator\n"), "announce name")
		Ω(co).Should(ContainSubstring("-help\n"), "flag prompts")
		Ω(co).Should(ContainSubstring("\nUsage of "), "Usage prefix")
	})

	It("executes with usage but no help on bad flag", func() {
		co, err := runOutput("-unknownFlag")
		Ω(err).Should(HaveOccurred(), jbmcExec)
		Ω(co).ShouldNot(ContainSubstring("\njava-buildpack-memory-calculator\n"), "announce name")
		Ω(co).Should(ContainSubstring("flag provided but not defined: "), "flag prompts")
		Ω(co).Should(ContainSubstring("-help\n"), "flag prompts")
		Ω(co).Should(ContainSubstring("\nUsage of "), "Usage prefix")
	})

	It("executes with error on unexpected argument", func() {
		unexpectedArgument := []string{"-totMemory=128m", "-spanishInquisition=surprise"}
		so, se, err := runOutAndErr(unexpectedArgument...)
		Ω(err).Should(HaveOccurred(), unexpectedArgument[1])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+unexpectedArgument[1])
		Ω(se).Should(ContainSubstring("flag provided but not defined: -spanishInquisition"), "stderr incorrect for "+unexpectedArgument[1])
	})

	It("executes with error when no total memory is supplied", func() {
		badFlags := []string{"-stackThreads=50", "-loadedClasses=100", "-vmOptions=-verbose:gc"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), "no -totMemory flag")

		Ω(string(so)).Should(BeEmpty(), "stdout not empty when no -totMemory flag")
		Ω(string(se)).Should(ContainSubstring("-totMemory must be specified"), "stderr incorrect when no -totMemory flag")
	})

	It("executes with error on bad total memory syntax", func() {
		badFlags := []string{"-totMemory=badmem"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Error in -totMemory flag: "), "stderr incorrect for "+badFlags[0])
	})

	It("executes with error when too little total memory is supplied", func() {
		badFlags := []string{"-totMemory=1023b"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Total memory (-totMemory flag) is less than 1K"), "stderr incorrect for "+badFlags[0])
	})

	It("executes with error when stackThreads is not supplied", func() {
		badFlags := []string{"-totMemory=128m", "-loadedClasses=100"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), "no -stackThreads flag")

		Ω(string(so)).Should(BeEmpty(), "stdout not empty when no -stackThreads flag")
		Ω(string(se)).Should(ContainSubstring("-stackThreads must be specified"), "stderr incorrect when no -stackThreads flag")
	})

	It("executes with error when stackThreads is negative", func() {
		badFlags :=
			[]string{"-totMemory=2G", "-stackThreads=-1"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Error in -stackThreads flag; value must be positive"), "stderr incorrect for "+badFlags[0])
	})

	It("executes with error when loadedClasses is not supplied", func() {
		badFlags := []string{"-totMemory=128m", "-stackThreads=10"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), "no -loadedClasses flag")

		Ω(string(so)).Should(BeEmpty(), "stdout not empty when no -loadedClasses flag")
		Ω(string(se)).Should(ContainSubstring("-loadedClasses must be specified"), "stderr incorrect when no -loadedClasses flag")
	})

	It("executes with error when loadedClasses is negative", func() {
		badFlags :=
			[]string{"-totMemory=2G", "-stackThreads=10", "-loadedClasses=-1"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Error in -loadedClasses flag; value must be positive"), "stderr incorrect for "+badFlags[0])
	})

	Context("with valid parameters", func() {
		var (
			totMemFlag string
			sOut, sErr []byte
			cmdErr     error
		)

		JustBeforeEach(func() {
			sOut, sErr, cmdErr = runOutAndErr(totMemFlag, "-stackThreads=10", "-loadedClasses=100")
		})

		Context("when there is sufficient total memory", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
			})

			It("succeeds", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal(""), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-XX:ReservedCodeCacheSize=240M",
					"-XX:CompressedClassSpaceSize=800K",
					"-Xmx3919899K",
					"-XX:MaxMetaspaceSize=7363K",
					"-XX:MaxDirectMemorySize=10M",
				), "stdout")
			})
		})

		Context("when there is insufficient total memory", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=32m"
			})

			It("fails with an error", func() {
				Ω(cmdErr).Should(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(ContainSubstring("Cannot calculate memory: insufficient memory remaining for heap (memory limit 32M < allocated memory 274404K)"),
					"stderr")
				Ω(string(sOut)).Should(Equal(""), "stdout")
			})
		})
	})
})

func runOutput(args ...string) ([]byte, error) {
	cmd := exec.Command(jbmcExec, args...)
	co, err := cmd.CombinedOutput()
	return co, err
}

func runOutAndErr(args ...string) ([]byte, []byte, error) {
	cmd := exec.Command(jbmcExec, args...)
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}
