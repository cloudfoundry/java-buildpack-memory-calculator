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
		Ω(co).Should(ContainSubstring("-help=false"), "flag prompts")
		Ω(co).Should(ContainSubstring("\nUsage of "), "Usage prefix")
	})

	It("executes with usage but no help on bad flag", func() {
		co, err := runOutput("-unknownFlag")
		Ω(err).Should(HaveOccurred(), jbmcExec)
		Ω(co).ShouldNot(ContainSubstring("\njava-buildpack-memory-calculator\n"), "announce name")
		Ω(co).Should(ContainSubstring("flag provided but not defined: "), "flag prompts")
		Ω(co).Should(ContainSubstring("-help=false"), "flag prompts")
		Ω(co).Should(ContainSubstring("\nUsage of "), "Usage prefix")
	})

	It("executes with error on bad total memory syntax", func() {
		badFlags := []string{"-totMemory=badmem"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Error in -totMemory flag: "), "stderr incorrect for "+badFlags[0])
	})

	It("executes with error on bad weights map", func() {
		badFlags := []string{"-memoryWeights=heap:-2", "-totMemory=4g"}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Weight must be positive in -memoryWeights flag; clause 'heap:-2'"), "stderr incorrect for "+badFlags[0])
	})

	Context("with valid parameters", func() {
		var (
			totMemFlag, weightsFlag, sizesFlag string
			sOut, sErr                         []byte
			cmdErr                             error
		)

		BeforeEach(func() {
			totMemFlag = "-totMemory=4g"
			weightsFlag = "-memoryWeights=heap:5,stack:1,permgen:3,native:1"
			sizesFlag = "-memorySizes=stack:1m"
		})

		JustBeforeEach(func() {
			goodFlags := []string{totMemFlag, weightsFlag, sizesFlag}
			sOut, sErr, cmdErr = runOutAndErr(goodFlags...)
		})

		Context("when no total memory is supplied", func() {
			BeforeEach(func() {
				totMemFlag = ""
				sizesFlag = "-memorySizes=stack:2m..,heap:30m..400m,permgen:10m..12m"
			})

			It("allocates the minima", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal(""), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-Xmx30M",
					"-Xms30M",
					"-Xss2M",
					"-XX:MaxPermSize=10M",
					"-XX:PermSize=10M",
				), "stdout")
			})
		})

		Context("using nothing but total memory parameter", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = ""
				sizesFlag = ""
			})

			It("succeeds", func() {
				// Ω(string(sErr)).Should(Equal(""), "stderr") // actually get a warning!
				Ω(string(sOut)).Should(Equal(""), "stdout")
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
			})
		})

		Context("using permgen", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = "-memoryWeights=heap:5,stack:1,permgen:3,native:1"
			})

			It("succeeds", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal(""), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-Xmx2G",
					"-Xms2G",
					"-Xss1M",
					"-XX:MaxPermSize=1258291K",
					"-XX:PermSize=1258291K",
				), "stdout")
			})
		})

		Context("using metaspace", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = "-memoryWeights=heap:5,stack:1,metaspace:3,native:1"
			})

			It("succeeds", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal(""), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-Xmx2G",
					"-Xms2G",
					"-Xss1M",
					"-XX:MaxMetaspaceSize=1258291K",
					"-XX:MetaspaceSize=1258291K",
				), "stdout")
			})
		})

		Context("when the specified maximum memory is larger than the total memory size", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = "-memoryWeights=heap:5,stack:1,permgen:3,native:1"
				sizesFlag = "-memorySizes=heap:3g,permgen:2g"
			})

			It("fails with an error", func() {
				Ω(cmdErr).Should(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(ContainSubstring("Cannot balance memory: Memory allocation failed for configuration: ["), "stderr")
				Ω(string(sOut)).Should(Equal(""), "stdout")
			})
		})

		Context("when the specified maximum memory sizes imply the total memory size may be too large", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = "-memoryWeights=heap:5,stack:1,permgen:3,native:1"
				sizesFlag = "-memorySizes=heap:800m,permgen:800m"
			})

			It("issues a warning", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal("There is more than 3 times more spare native memory than the default so configured Java memory may be too small or available memory may be too large\n"), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-Xmx800M",
					"-Xms800M",
					"-Xss3120K",
					"-XX:MaxPermSize=800M",
					"-XX:PermSize=800M",
				), "stdout")
			})
		})

		Context("when the specified maximum memory sizes, including native, imply the total memory size may be too large", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = "-memoryWeights=heap:5,stack:1,permgen:3,native:1"
				sizesFlag = "-memorySizes=heap:1m,permgen:1m,stack:2m,native:2000m"
			})

			It("issues a warning", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal("The allocated Java memory sizes total 2469478K which is less than 0.8 of the available memory, so configured Java memory sizes may be too small or available memory may be too large\n"), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-XX:MaxPermSize=1M",
					"-XX:PermSize=1M",
					"-Xmx1M",
					"-Xms1M",
					"-Xss2M",
				), "stdout")
			})
		})

		Context("when the specified maximum heap size is close to the default", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = "-memoryWeights=heap:5,stack:1,permgen:3,native:1"
				sizesFlag = "-memorySizes=heap:2049m"
			})

			It("issues a warning", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal("The specified value 2049M for memory type heap is close to the computed value 2G. Consider taking the default.\n"), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-XX:MaxPermSize=1257676K",
					"-XX:PermSize=1257676K",
					"-Xmx2049M",
					"-Xms2049M",
					"-Xss1023K",
				), "stdout")
			})
		})

		Context("when the specified maximum permgen size is close to the default", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = "-memoryWeights=heap:5,stack:1,permgen:3,native:1"
				sizesFlag = "-memorySizes=permgen:1339m"
			})

			It("issues a warning", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal("The specified value 1339M for memory type permgen is close to the computed value 1258291K. Consider taking the default.\n"), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-Xss984K",
					"-XX:MaxPermSize=1339M",
					"-XX:PermSize=1339M",
					"-Xmx2016548K",
					"-Xms2016548K",
				), "stdout")
			})
		})
	})
})

func runOutput(args ...string) ([]byte, error) {
	cmd := exec.Command(jbmcExec, squash(args)...)
	co, err := cmd.CombinedOutput()
	return co, err
}

func runOutAndErr(args ...string) ([]byte, []byte, error) {
	cmd := exec.Command(jbmcExec, squash(args)...)
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout, cmd.Stderr = stdout, stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func squash(ss []string) []string {
	result := []string{}
	for _, s := range ss {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
