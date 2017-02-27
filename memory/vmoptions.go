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
package memory

import (
	"bytes"
	"strings"
)

//go:generate counterfeiter -o vmoptionsfakes/fake_vmoptions.go . VmOptions
type VmOptions interface {
	DeltaString() string
	MemOpt(memoryType MemoryType) MemSize
	SetMemOpt(memoryType MemoryType, size MemSize)
}

type vmOptions struct {
	rawOpts      string
	memOptWasRaw map[MemoryType]bool
	memOpts      map[MemoryType]MemSize
}

type MemoryType int

const (
	MaxHeapSize MemoryType = iota
	MaxMetaspaceSize
	StackSize
	MaxDirectMemorySize
	ReservedCodeCacheSize
	CompressedClassSpaceSize
	MaxPermSize
)

var switches = map[MemoryType]string{
	MaxHeapSize:              "-Xmx",
	MaxMetaspaceSize:         "-XX:MaxMetaspaceSize=",
	StackSize:                "-Xss",
	MaxDirectMemorySize:      "-XX:MaxDirectMemorySize=",
	ReservedCodeCacheSize:    "-XX:ReservedCodeCacheSize=",
	CompressedClassSpaceSize: "-XX:CompressedClassSpaceSize=",
	MaxPermSize:              "-XX:MaxPermSize=",
}

func NewVmOptions(rawOpts string) (*vmOptions, error) {
	var mo map[MemoryType]MemSize = map[MemoryType]MemSize{}
	var mowr map[MemoryType]bool = map[MemoryType]bool{}

	for optMemoryType, sw := range switches {
		opt, err := parseOpt(rawOpts, sw)
		if err != nil {
			return nil, err
		}
		if opt != MEMSIZE_ZERO {
			mo[optMemoryType] = opt
			mowr[optMemoryType] = true
		} else {
			mowr[optMemoryType] = false
		}
	}

	return &vmOptions{
		rawOpts:      rawOpts,
		memOptWasRaw: mowr,
		memOpts:      mo,
	}, nil
}

func (vm *vmOptions) DeltaString() string {
	var bb bytes.Buffer

	first := true
	for k, v := range vm.memOpts {
		if vm.memOptWasRaw[k] {
			continue
		}
		if !first {
			bb.WriteString(" ")
		}
		bb.WriteString(switches[k] + v.String())
		first = false
	}

	return bb.String()
}

func (vm *vmOptions) MemOpt(memoryType MemoryType) MemSize {
	return vm.memOpts[memoryType]
}

func (vm *vmOptions) SetMemOpt(memoryType MemoryType, size MemSize) {
	vm.memOpts[memoryType] = size
}

func parseOpt(rawOpts string, sw string) (MemSize, error) {
	opts := strings.Split(rawOpts, " ")
	for _, opt := range opts {
		if opt == "" {
			continue
		}
		if strings.Index(opt, sw) == 0 {
			value := opt[len(sw):]
			return NewMemSizeFromString(value)
		}
	}
	return MEMSIZE_ZERO, nil
}
