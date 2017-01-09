// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2015-2016 the original author or authors.
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

package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
)

const (
	executableName    = "java-buildpack-memory-calculator"
	totalFlag         = "totMemory"
	threadsFlag       = "stackThreads"
	loadedClassesFlag = "loadedClasses"
	vmOptionsFlag     = "vmOptions"
)

func printHelp() {
	fmt.Printf("\n%s\n", executableName)
	fmt.Println("\nCalculate JRE memory switches " +
		"based upon the total memory available and the size ranges and weights given.\n")
	flag.Usage()
}

var (
	help = flag.Bool("help", false, "prints description and flag help")

	totMemory = flag.String(totalFlag, "",
		"total memory available to allocate, expressed as an integral "+
			"number of bytes (B), kilobytes (K), megabytes (M) or gigabytes (G), e.g. '1G'")
	stackThreads = flag.Int(threadsFlag, 0,
		"number of threads to use in stack allocation calculations'")
	loadedClasses = flag.Int(loadedClassesFlag, 0,
		"an estimate of the number of classes which will be loaded when the application is running")
	vmOptions = flag.String(vmOptionsFlag, "",
		"Java VM options, typically the JAVA_OPTS specified by the user")
)

// Validate flags passed on command line; exit(1) if invalid; exit(2) if help printed
func ValidateFlags() (memSize memory.MemSize, numThreads int, numLoadedClasses int, vmOpts string) {

	flag.Parse() // exit on error

	if noArgs(os.Args[1:]) || helpArg() {
		printHelp()
		os.Exit(2)
	}

	// validation routines will not return on error
	validateNoArguments()
	memSize = validateTotMemory(*totMemory)
	numThreads = validateNumThreads(*stackThreads)
	numLoadedClasses = validateLoadedClasses(*loadedClasses)
	vmOpts = *vmOptions

	return
}

func validateNoArguments() {
	if len(flag.Args()) != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument: %s\n", flag.Args()[0])
		os.Exit(1)
	}
}

func validateTotMemory(mem string) memory.MemSize {
	if mem == "" {
		fmt.Fprintf(os.Stderr, "-%s must be specified", totalFlag)
		os.Exit(1)
	}
	ms, err := memory.NewMemSizeFromString(mem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in -%s flag: %s", totalFlag, err)
		os.Exit(1)
	}
	if ms.LessThan(memory.MemSize(1024)) {
		fmt.Fprintf(os.Stderr, "Total memory (-%s flag) is less than 1K", totalFlag)
		os.Exit(1)
	}
	return ms
}

func validateLoadedClasses(loadedClasses int) int {
	if loadedClasses == 0 {
		fmt.Fprintf(os.Stderr, "-%s must be specified", loadedClassesFlag)
		os.Exit(1)
	}
	if loadedClasses < 0 {
		fmt.Fprintf(os.Stderr, "Error in -%s flag; value must be positive", loadedClassesFlag)
		os.Exit(1)
	}
	return loadedClasses
}

func validateNumThreads(stackThreads int) int {
	if stackThreads == 0 {
		fmt.Fprintf(os.Stderr, "-%s must be specified", threadsFlag)
		os.Exit(1)
	}
	if stackThreads < 0 {
		fmt.Fprintf(os.Stderr, "Error in -%s flag; value must be positive", threadsFlag)
		os.Exit(1)
	}
	return stackThreads
}

func noArgs(args []string) bool {
	return len(args) == 0
}

func helpArg() bool {
	return flag.Parsed() && *help
}
