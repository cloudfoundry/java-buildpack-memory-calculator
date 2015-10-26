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

	It("executes with error when no total memory is supplied", func() {
		badFlags :=
			[]string{
				"-memoryWeights=heap:5,stack:1,permgen:3,native:1",
				"-memorySizes=stack:2m..,heap:30m..400m,permgen:10m..12m",
			}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), "no -totMemory flag")

		Ω(string(so)).Should(BeEmpty(), "stdout not empty when no -totMemory flag")
		Ω(string(se)).Should(ContainSubstring("-totMemory must be specified"), "stderr incorrect when no -totMemory flag")
	})

	It("executes with error when too little total memory is supplied", func() {
		badFlags :=
			[]string{
				"-totMemory=1023b",
				"-memoryWeights=heap:5,stack:1,permgen:3,native:1",
				"-memorySizes=stack:2m..,heap:30m..400m,permgen:10m..12m",
			}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Total memory (-totMemory flag) is less than 1K"), "stderr incorrect for "+badFlags[0])
	})

	It("executes with error when initial is not a percentage", func() {
		badFlags :=
			[]string{
				"-totMemory=2G",
				"-memoryWeights=heap:5,stack:1,permgen:3,native:1",
				"-memorySizes=stack:2m..,heap:30m..400m,permgen:10m..12m",
				"-memoryInitials=heap:50",
			}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Bad initial value in -memoryInitials flag; clause 'heap:50' : value must be a percentage (e.g. 10%)"), "stderr incorrect for "+badFlags[0])
	})

	It("executes with error when initial too big", func() {
		badFlags :=
			[]string{
				"-totMemory=2G",
				"-memoryWeights=heap:5,stack:1,permgen:3,native:1",
				"-memorySizes=stack:2m..,heap:30m..400m,permgen:10m..12m",
				"-memoryInitials=heap:101%",
			}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Initial value must be zero or more but no more than 100% in -memoryInitials flag; clause 'heap:101%'"), "stderr incorrect for "+badFlags[0])
	})

	It("executes with error when initial is negative", func() {
		badFlags :=
			[]string{
				"-totMemory=2G",
				"-memoryWeights=heap:5,stack:1,permgen:3,native:1",
				"-memorySizes=stack:2m..,heap:30m..400m,permgen:10m..12m",
				"-memoryInitials=heap:-1%",
			}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Initial value must be zero or more but no more than 100% in -memoryInitials flag; clause 'heap:-1%'"), "stderr incorrect for "+badFlags[0])
	})

	It("executes with error when stackThreads is negative", func() {
		badFlags :=
			[]string{
				"-totMemory=2G",
				"-memoryWeights=heap:5,stack:1,permgen:3,native:1",
				"-memorySizes=stack:2m..,heap:30m..400m,permgen:10m..12m",
				"-stackThreads=-1",
			}
		so, se, err := runOutAndErr(badFlags...)
		Ω(err).Should(HaveOccurred(), badFlags[0])

		Ω(string(so)).Should(BeEmpty(), "stdout not empty for "+badFlags[0])
		Ω(string(se)).Should(ContainSubstring("Error in -stackThreads flag; value must be positive"), "stderr incorrect for "+badFlags[0])
	})

	Context("with valid parameters", func() {
		var (
			totMemFlag, weightsFlag, sizesFlag, initialsFlag string
			sOut, sErr                                       []byte
			cmdErr                                           error
		)

		BeforeEach(func() {
			totMemFlag = "-totMemory=4g"
			weightsFlag = "-memoryWeights=heap:5,stack:1,permgen:3,native:1"
			sizesFlag = "-memorySizes=stack:1m"
			initialsFlag = "-memoryInitials=heap:50%,permgen:50%"
		})

		JustBeforeEach(func() {
			goodFlags := []string{totMemFlag, weightsFlag, sizesFlag, initialsFlag}
			sOut, sErr, cmdErr = runOutAndErr(goodFlags...)
		})

		Context("using nothing but total memory parameter", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = ""
				sizesFlag = ""
				initialsFlag = ""
			})

			It("succeeds", func() {
				// Ω(string(sErr)).Should(Equal(""), "stderr") // actually get a warning!
				Ω(string(sOut)).Should(Equal(""), "stdout")
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
			})
		})

		Context("using no memoryInitials parameter", func() {
			BeforeEach(func() {
				initialsFlag = ""
			})

			It("succeeds", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal(""), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-Xmx2G",
					"-Xss1M",
					"-XX:MaxPermSize=1258291K",
				), "stdout")
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
					"-Xms1G",
					"-Xss1M",
					"-XX:MaxPermSize=1258291K",
					"-XX:PermSize=629145K",
				), "stdout")
			})
		})

		Context("using metaspace", func() {
			BeforeEach(func() {
				totMemFlag = "-totMemory=4g"
				weightsFlag = "-memoryWeights=heap:5,stack:1,metaspace:3,native:1"
				initialsFlag = "-memoryInitials=heap:50%,metaspace:50%"
			})

			It("succeeds", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal(""), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-Xmx2G",
					"-Xms1G",
					"-Xss1M",
					"-XX:MaxMetaspaceSize=1258291K",
					"-XX:MetaspaceSize=629145K",
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
					"-Xms400M",
					"-Xss3120K",
					"-XX:MaxPermSize=800M",
					"-XX:PermSize=400M",
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
				Ω(string(sErr)).Should(ContainSubstring("The allocated Java memory sizes total 2469478K which is less than 0.8 of the available memory, so configured Java memory sizes may be too small or available memory may be too large\n"), "stderr")
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
					"-XX:PermSize=628838K",
					"-Xmx2049M",
					"-Xms1049088K",
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
					"-XX:PermSize=685568K",
					"-Xmx2016548K",
					"-Xms1008274K",
				), "stdout")
			})
		})
		Context("when the specified initial memory is less than static minimum", func() {
			BeforeEach(func() {
				initialsFlag = "-memoryInitials=heap:0%"
			})

			It("issues a warning", func() {
				Ω(cmdErr).ShouldNot(HaveOccurred(), "exit status")
				Ω(string(sErr)).Should(Equal("The configured initial memory size 0 for heap is less than the minimum 2M.  Setting initial value to 2M.\n"), "stderr")
				Ω(strings.Split(string(sOut), " ")).Should(ConsistOf(
					"-Xss1M",
					"-XX:MaxPermSize=1258291K",
					"-Xmx2G",
					"-Xms2M",
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
