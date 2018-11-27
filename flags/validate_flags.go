// Encoding: utf-8
// Cloud Foundry Java Buildpack
// Copyright (c) 2015-2018 the original author or authors.
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
	"strconv"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
	sysmem "github.com/pbnjay/memory"
)

const (
	executableName    = "java-buildpack-memory-calculator"
	totalFlag         = "totMemory"
	threadsFlag       = "stackThreads"
	loadedClassesFlag = "loadedClasses"
	vmOptionsFlag     = "vmOptions"
	poolTypeFlag      = "poolType"
	headRoomFlag      = "headRoom"
)

func printHelp() {
	fmt.Printf("\n%s\n", executableName)
	fmt.Printf("\nCalculates JVM memory switches based on the total memory available, the number of classes the application will load, "+
		"the number of threads that will be used, and any JVM options provided as input.\n\n"+
		"The output consists of any calculated memory switches.\n\n"+
		"If a calculated memory switch value is unsuitable, it can be set in the JVM options provided as input and will no longer be calculated.\n\n"+
		"Example invocation from a shell:\n"+
		"$ %s -loadedClasses=1000 -stackThreads=10 -totMemory=1g -poolType=metaspace -vmOptions=-XX:MaxDirectMemorySize=100M\\ -verbose:gc\n\n", executableName)
	flag.Usage()
}

var (
	help = flag.Bool("help", false, "prints description and flag help")

	totMemory = flag.String(totalFlag, "",
		"total memory available to allocate, expressed as an integral "+
			"number of bytes (B), kilobytes (K), megabytes (M) or gigabytes (G), e.g. '1G'")
	stackThreads = flag.Int(threadsFlag, 0,
		"number of threads that will be used")
	loadedClasses = flag.Int(loadedClassesFlag, 0,
		"an estimate of the number of classes that will be loaded when the application is running")
	vmOptions = flag.String(vmOptionsFlag, "",
		"Java VM options, typically the JAVA_OPTS specified by the user")
	poolType = flag.String(poolTypeFlag, "",
		"the type of JVM pool used in the calculation. Set this to 'permgen' for Java 7 and to 'metaspace' for Java 8 and later.")
	headRoom = flag.Float64(headRoomFlag, 0,
		"percentage of total memory available which will be left unallocated to cover JVM overheads")
)

// Validate flags passed on command line; exit(1) if invalid; exit(2) if help printed
func ValidateFlags() (memSize memory.MemSize, numThreads int, numLoadedClasses int, pType string, vmOpts string, hdRoom float64) {

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
	pType = validatePoolType(*poolType)
	vmOpts = *vmOptions
	hdRoom = validateHeadRoom(*headRoom)

	return
}

func validatePoolType(poolType string) string {
	if poolType == "" {
		fmt.Fprintf(os.Stderr, "-%s must be specified", poolTypeFlag)
		os.Exit(1)
	}
	if poolType != "permgen" && poolType != "metaspace" {
		fmt.Fprintf(os.Stderr, "Error in -%s flag: must be 'permgen' or 'metaspace'", poolTypeFlag)
		os.Exit(1)
	}
	return poolType
}

func validateNoArguments() {
	if len(flag.Args()) != 0 {
		fmt.Fprintf(os.Stderr, "unexpected argument: %s\n", flag.Args()[0])
		os.Exit(1)
	}
}

func validateTotMemory(mem string) memory.MemSize {
	if mem == "" {
		mem = strconv.FormatUint(sysmem.TotalMemory(),10) + "b"
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

func validateHeadRoom(headRoom float64) float64 {
	if headRoom < 0 || headRoom > 100 {
		fmt.Fprintf(os.Stderr, "Head room (-%s) is not a valid percentage: %f", headRoomFlag, headRoom)
		os.Exit(1)
	}
	return headRoom
}

func noArgs(args []string) bool {
	return len(args) == 0
}

func helpArg() bool {
	return flag.Parsed() && *help
}
