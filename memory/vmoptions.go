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
package memory

import (
	"bytes"
	"strings"
)

//go:generate counterfeiter -o vmoptionsfakes/fake_vmoptions.go . VmOptions
type VmOptions interface {
	DeltaString() string
	Copy() VmOptions
	String() string
	MemOpt(memoryType MemoryType) MemSize
	SetMemOpt(memoryType MemoryType, size MemSize)
	ClearMemOpt(memoryType MemoryType) // forgets a memory type, raw or not
}

type vmOptions struct {
	memOpts      map[MemoryType]MemSize
	memOptWasRaw map[MemoryType]bool
}

type MemoryType int

const (
	MaxHeapSize MemoryType = iota
	StackSize
	ReservedCodeCacheSize
	MaxDirectMemorySize
	MaxMetaspaceSize
	MaxPermSize
	MemoryTypeLimit // not an actual memory type, used for enumerating the valid memory types
)

var switches = map[MemoryType]string{
	MaxHeapSize:              "-Xmx",
	MaxMetaspaceSize:         "-XX:MaxMetaspaceSize=",
	StackSize:                "-Xss",
	MaxDirectMemorySize:      "-XX:MaxDirectMemorySize=",
	ReservedCodeCacheSize:    "-XX:ReservedCodeCacheSize=",
	MaxPermSize:              "-XX:MaxPermSize=",
}

func NewVmOptions(rawOpts string) (*vmOptions, error) {
	var mo map[MemoryType]MemSize = map[MemoryType]MemSize{}
	var mowr map[MemoryType]bool = map[MemoryType]bool{}

	for optMemoryType, sw := range switches {
		opt, raw, err := parseOpt(rawOpts, sw)
		if err != nil {
			return nil, err
		}
		if raw {
			mo[optMemoryType] = opt
			mowr[optMemoryType] = true
		} else {
			mowr[optMemoryType] = false
		}
	}

	return &vmOptions{
		memOptWasRaw: mowr,
		memOpts:      mo,
	}, nil
}

func (vm *vmOptions) Copy() VmOptions {
	copy := &vmOptions{}

	copy.memOpts = map[MemoryType]MemSize{}
	for k, v := range vm.memOpts {
		copy.memOpts[k] = v
	}

	copy.memOptWasRaw = map[MemoryType]bool{}
	for k, v := range vm.memOptWasRaw {
		copy.memOptWasRaw[k] = v
	}

	return copy
}

func (vm *vmOptions) DeltaString() string {
	var bb bytes.Buffer

	first := true
	for k := MemoryType(0); k < MemoryTypeLimit; k++ {
		v, ok := vm.memOpts[k]
		if !ok {
			continue
		}
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

// Human readable form of the VM options with a fixed order.
func (vm *vmOptions) String() string {
	var bb bytes.Buffer

	first := true
	for k := MemoryType(0); k < MemoryTypeLimit; k++ {
		v, ok := vm.memOpts[k]
		if !ok {
			continue
		}
		if !first {
			bb.WriteString(", ")
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

func (vm *vmOptions) ClearMemOpt(memoryType MemoryType) {
	delete(vm.memOpts, memoryType)
	delete(vm.memOptWasRaw, memoryType)
}

func parseOpt(rawOpts string, sw string) (MemSize, bool, error) {
	optValue := MEMSIZE_ZERO
	raw := false
	opts := strings.Split(rawOpts, " ")
	for _, opt := range opts {
		if opt == "" {
			continue
		}
		if strings.Index(opt, sw) == 0 {
			raw = true
			value := opt[len(sw):]
			var err error
			optValue, err = NewMemSizeFromString(value)
			if err != nil {
				return optValue, raw, err
			}
		}
	}
	return optValue, raw, nil
}
