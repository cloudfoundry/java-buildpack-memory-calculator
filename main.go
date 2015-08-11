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

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/flags"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory/switches"
)

const (
	exec_name = "java-buildpack-memory-calculator"
)

func main() {

	// validateFlags() will exit on error
	memSize, numThreads, weights, sizes, initials := flags.ValidateFlags()

	allocator, err := memory.NewAllocator(sizes, weights)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot allocate memory: %s", err)
		os.Exit(1)
	}

	if err = allocator.Balance(memSize, numThreads); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot balance memory: %s", err)
		os.Exit(1)
	}

	allocator.GenerateInitialAllocations(initials)

	allocatorSwitches := allocator.Switches(switches.AllocatorJreSwitchFuns)

	if warnings := allocator.GetWarnings(); len(warnings) != 0 {
		fmt.Fprintln(os.Stderr, strings.Join(warnings, "\n"))
	}

	fmt.Fprint(os.Stdout, strings.Join(allocatorSwitches, " "))

}
