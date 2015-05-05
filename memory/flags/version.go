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
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Micro int
}

func NewVersion(vstr string) (Version, error) {
	subVersions := strings.Split(vstr, ".")

	if ln := len(subVersions); ln == 0 || ln > 3 {
		return Version{}, fmt.Errorf("Version ('%s') is empty, or has too many components.", vstr)
	}

	var err error
	v := [3]int{0, 0, 0}

	for i, sv := range subVersions {
		v[i], err = strconv.Atoi(sv)
		if err != nil || v[i] < 0 {
			return Version{}, fmt.Errorf("Version ('%s') has incorrect format", vstr)
		}
	}
	return Version{v[0], v[1], v[2]}, nil
}

func (v Version) LessThan(o Version) bool {
	if vMaj, oMaj := v.Major, o.Major; vMaj < oMaj {
		return true
	} else if vMaj == oMaj {
		if vMin, oMin := v.Minor, o.Minor; vMin < oMin {
			return true
		} else if vMin == oMin {
			if vMic, oMic := v.Micro, o.Micro; vMic < oMic {
				return true
			}
		}
	}
	return false
}
