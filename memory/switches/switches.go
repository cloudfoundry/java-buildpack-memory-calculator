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

package switches

import "fmt"

type Funs map[string]switchFun

type switchFun func(v string) []string

// Returns a function ƒa • (ss map ƒs•Sprintf(s,a))
// This takes a string v and returns a slice of strings, formatted with each string in ss, parameter v.
func apply(ss ...string) func(string) []string {
	return func(v string) []string {
		var res = make([]string, 0, 2)
		for _, form := range ss {
			res = append(res, fmt.Sprintf(form, v))
		}
		return res
	}
}

var JreSwitchFuns = Funs{
	"heap":      apply("-Xmx%s", "-Xms%s"),
	"metaspace": apply("-XX:MaxMetaspaceSize=%s", "-XX:MetaspaceSize=%s"),
	"permgen":   apply("-XX:MaxPermSize=%s", "-XX:PermSize=%s"),
	"stack":     apply("-Xss%s"),
}

func (sf Funs) HasKey(akey string) bool {
	_, ok := sf[akey]
	return ok
}

func (sf Funs) Apply(akey string, aparm string) []string {
	if ƒ, ok := sf[akey]; ok {
		return ƒ(aparm)
	}
	return make([]string, 0, 0)
}
