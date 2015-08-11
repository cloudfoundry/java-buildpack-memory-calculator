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

package flags

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
)

const (
	executableName = "java-buildpack-memory-calculator"
	totalFlag      = "totMemory"
	weightsFlag    = "memoryWeights"
	sizesFlag      = "memorySizes"
	initialsFlag   = "memoryInitials"
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
			"number of bytes, kilobytes, megabytes or gigabytes")
	memoryWeights = flag.String(weightsFlag, "",
		"the weights given to each memory type, e.g. 'heap:15,permgen:5,stack:1,native:2'")
	memorySizes = flag.String(sizesFlag, "",
		"the size ranges allowed for each memory type, "+
			"e.g. 'heap:128m..1G,permgen:64m,stack:2m..4m,native:100m..'")
	memoryInitials = flag.String(initialsFlag, "",
		"the initial values for each memory type, "+
			"e.g. 'heap:128m,permgen:64m'")
)

// Validate flags passed on command line; exit(1) if invalid; exit(2) if help printed
func ValidateFlags() (memSize memory.MemSize, weights map[string]float64, sizes map[string]memory.Range, initials map[string]float64) {

	flag.Parse() // exit on error

	if noArgs(os.Args[1:]) || helpArg() {
		printHelp()
		os.Exit(2)
	}

	// validation routines will not return on error
	memSize = validateTotMemory(*totMemory)
	weights = validateWeights(*memoryWeights)
	sizes = validateSizes(*memorySizes)
	initials = validateInitials(*memoryInitials)

	return memSize, weights, sizes, initials
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

func validateWeights(weights string) map[string]float64 {
	ws := map[string]float64{}

	if weights == "" {
		return ws
	}

	weightClauses := strings.Split(weights, ",")
	for _, clause := range weightClauses {
		if parts := strings.Split(clause, ":"); len(parts) == 2 {
			if floatVal, err := strconv.ParseFloat(parts[1], 32); err != nil {
				fmt.Fprintf(os.Stderr, "Bad weight in -%s flag; clause '%s' : %s", weightsFlag, clause, err)
				os.Exit(1)
			} else if floatVal <= 0.0 {
				fmt.Fprintf(os.Stderr, "Weight must be positive in -%s flag; clause '%s'", weightsFlag, clause)
				os.Exit(1)
			} else {
				ws[parts[0]] = floatVal
			}
		} else {
			fmt.Fprintf(os.Stderr, "Bad clause '%s' in -%s flag", clause, weightsFlag)
			os.Exit(1)
		}
	}

	return ws
}

func validateSizes(sizes string) map[string]memory.Range {
	rs := map[string]memory.Range{}

	if sizes == "" {
		return rs
	}

	rangeClauses := strings.Split(sizes, ",")
	for _, clause := range rangeClauses {
		if parts := strings.Split(clause, ":"); len(parts) == 2 {
			var err error
			if rs[parts[0]], err = memory.NewRangeFromString(parts[1]); err != nil {
				fmt.Fprintf(os.Stderr, "Bad range in -%s flag, clause '%s' : %s", sizesFlag, clause, err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Bad clause '%s' in -%s flag", clause, sizesFlag)
			os.Exit(1)
		}
	}

	return rs
}

func validateInitials(initials string) map[string]float64 {
	is := map[string]float64{}

	if initials == "" {
		return is
	}

	initialClauses := strings.Split(initials, ",")
	for _, clause := range initialClauses {
		if parts := strings.Split(clause, ":"); len(parts) == 2 {
			if !strings.HasSuffix(parts[1], "%") {
				fmt.Fprintf(os.Stderr, "Bad initial value in -%s flag; clause '%s' : value must be a percentage (e.g. 10%%)", initialsFlag, clause)
				os.Exit(1)
			}
			if floatVal, err := strconv.ParseFloat(strings.Replace(parts[1], "%", "", 1), 32); err != nil {
				fmt.Fprintf(os.Stderr, "Bad initial value in -%s flag; clause '%s' : %s", initialsFlag, clause, err)
				os.Exit(1)
			} else if floatVal < 0.0 || floatVal > 100.0 {
				fmt.Fprintf(os.Stderr, "Initial value must be zero or more but no more than 100%% in -%s flag; clause '%s'", initialsFlag, clause)
				os.Exit(1)
			} else {
				//Convert value to valid scale factor
				is[parts[0]] = floatVal * .01
			}
		} else {
			fmt.Fprintf(os.Stderr, "Bad clause '%s' in -%s flag", clause, initialsFlag)
			os.Exit(1)
		}
	}

	return is
}

func noArgs(args []string) bool {
	return len(args) == 0
}

func helpArg() bool {
	return flag.Parsed() && *help
}
