package flags

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/memory"
)

var (
	help = flag.Bool("help", false, "prints description and flag help")

	jreVersion = flag.String("jreVersion", "",
		"the version of Java runtime to use; "+
			"this determines the names and the format of the switches generated")
	totMemory = flag.String("totMemory", "",
		"total memory available to allocate, expressed as an integral "+
			"number of bytes, kilobytes, megabytes or gigabytes")
	memoryWeights = flag.String("memoryWeights", "",
		"the weights given to each memory type, e.g. 'heap:15,permgen:5,stack:1,native:2'")
	memorySizes = flag.String("memorySizes", "",
		"the size ranges allowed for each memory type, "+
			"e.g. 'heap:128m..1G,permgen:64m,stack:2m..4m,native:100m..'")
)

func ValidateFlags() (memory.MemSize, map[string]float64, map[string]memory.Range) {

	flag.Parse() // exit on error

	if noArgs(os.Args[1:]) || helpArg() {
		printHelp()
		os.Exit(2)
	}

	// validation routines exit on error
	version := validateJreVersion()
	memSize := validateTotMemory()
	weights := validateWeights(version)
	sizes := validateSizes(version)
	return memSize, weights, sizes

}

func validateJreVersion() Version {
	v, err := NewVersion(*jreVersion)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in -jreVersion: %s", err)
		os.Exit(1)
	}
	return v
}

func validateTotMemory() memory.MemSize {
	ms, err := memory.NewMemSizeFromString(*totMemory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in -totMemory: %s", err)
		os.Exit(1)
	}
	return ms
}

func validateWeights(version Version) map[string]float64 {
	ws := map[string]float64{}

	weightClauses := strings.Split(*memoryWeights, ",")
	for _, clause := range weightClauses {
		if parts := strings.Split(clause, ":"); len(parts) == 2 && validMemoryType(parts[0], version) {
			var err error
			if ws[parts[0]], err = strconv.ParseFloat(parts[1], 32); err != nil {
				fmt.Fprintf(os.Stderr, "Bad float in -memoryWeights, clause %s : %s", clause, err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Bad clause '%s' in -memoryWeights", clause)
			os.Exit(1)
		}
	}

	return ws
}

func validateSizes(version Version) map[string]memory.Range {
	rs := map[string]memory.Range{}

	rangeClauses := strings.Split(*memorySizes, ",")
	for _, clause := range rangeClauses {
		if parts := strings.Split(clause, ":"); len(parts) == 2 && validMemoryType(parts[0], version) {
			var err error
			if rs[parts[0]], err = memory.NewRangeFromString(parts[1]); err != nil {
				fmt.Fprintf(os.Stderr, "Bad range in -memorySizes, clause %s : %s", clause, err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Bad clause '%s' in -memorySizes", clause)
			os.Exit(1)
		}
	}

	return rs
}

func printHelp() {
	fmt.Printf("\n%s\n", "java-buildpack-memory-calculator")
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

func validMemoryType(memType string, version Version) bool {
	return true
}
