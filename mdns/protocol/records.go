// Copyright (C) 2022 The go-mdns Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	"fmt"
	"strconv"
)

// Records is a list of Record.
type Records []Record

// String returns the string representation.
func (records Records) String() string {
	type record []string

	lines := []record{}
	for _, r := range records {
		record := record{
			r.Name(),
			r.Type().String(),
			r.Content(),
		}
		lines = append(lines, record)
	}

	maxRecordLen := []int{0, 0, 0}
	for _, r := range lines {
		for n, s := range r {
			if maxRecordLen[n] < len(s) {
				maxRecordLen[n] = len(s)
			}
		}
	}

	str := ""
	for _, r := range lines {
		for n, s := range r {
			sfmt := "%" + strconv.Itoa(maxRecordLen[n]) + "s"
			str += fmt.Sprintf(sfmt, s)
			if n < len(r)-1 {
				str += " "
			}
		}
		str += "\n"
	}

	return str
}
