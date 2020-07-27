/*
 * Copyright 2015-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/java-buildpack-memory-calculator/v4/calculator"
	"github.com/cloudfoundry/java-buildpack-memory-calculator/v4/flags"
	flag "github.com/spf13/pflag"
)

func main() {
	h := flags.DefaultHeadRoom
	j := flags.DefaultJVMOptions
	l := flags.DefaultLoadedClassCount
	t := flags.DefaultThreadCount
	m := flags.DefaultTotalMemory

	c := calculator.Calculator{HeadRoom: &h, JvmOptions: &j, LoadedClassCount: &l, ThreadCount: &t, TotalMemory: &m}

	flag.Var(c.HeadRoom, flags.FlagHeadRoom, "percentage of total memory available which will be left unallocated to cover JVM overhead")
	flag.Var(c.JvmOptions, flags.FlagJVMOptions, "JVM options, typically JAVA_OPTS")
	flag.Var(c.LoadedClassCount, flags.FlagLoadedClassCount, "the number of classes that will be loaded when the application is running")
	flag.Var(c.ThreadCount, flags.FlagThreadCount, "the number of user threads")
	flag.Var(c.TotalMemory, "total-memory", "total memory available to the application, typically expressed with size classification (B, K, M, G, T)")
	flag.Parse()

	if !validate(c.HeadRoom, c.JvmOptions, c.LoadedClassCount, c.ThreadCount, c.TotalMemory) {
		_, _ = fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	o, err := c.Calculate()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}

	s := make([]string, len(o))
	for i, t := range o {
		s[i] = t.String()
	}

	fmt.Println(output(o))
}

func output(o []fmt.Stringer) string {
	s := make([]string, len(o))

	for i, t := range o {
		s[i] = t.String()
	}

	return strings.Join(s, " ")
}

func validate(vs ...flags.Validatable) bool {
	valid := true

	for _, v := range vs {
		if err := v.Validate(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			valid = false
		}
	}

	return valid
}
