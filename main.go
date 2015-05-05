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
	"flag"
	"fmt"
	"os"
)

const (
	exec_name = "java-buildpack-memory-calculator"
)

var (
	help = flag.Bool("help", false, "prints description and flag help")

	totMemory = flag.String("totMemory", "",
		"total memory available to allocate, expressed as an integral "+
			"number of bytes, kilobytes, megabytes or gigabytes")
	jreVersion = flag.String("jreVersion", "",
		"the version of Java runtime to use; "+
			"this determines the names and the format of the switches generated")
	memoryWeights = flag.String("memoryWeights", "",
		"the weights given to each memory type, e.g. 'heap:15,permgen:5,stack:1,native:2'")
	memorySizes = flag.String("memorySizes", "",
		"the size ranges allowed for each memory type, "+
			"e.g. 'heap:128m..1G,permgen:64m,stack:2m..4m,native:100m..'")
)

func main() {
	flag.Parse() // exit on error

	if noArgs(os.Args[1:]) || helpArg() {
		printHelp()
		os.Exit(2)
	}

	os.Exit(1) // error if we get here
}

func printHelp() {
	fmt.Printf("\n%s\n", exec_name)
	fmt.Println("\nCalculate JRE memory switches " +
		"based upon the total memory available and the size ranges and weights given.\n")
	flag.Usage()
}

func noArgs(args []string) bool {
	return len(args) == 0
}

func helpArg() bool {
	return flag.Parsed() && *help
}
