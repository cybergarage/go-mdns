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

package cmd

import (
	"fmt"
	"strings"
)

// Format represents the output format.
type Format int

// The output formats.
const (
	// FormatTable prints the results as an aligned table.
	FormatTable Format = iota
	// FormatJSON prints the results as JSON, one object a line.
	FormatJSON
	// FormatCSV prints the results as CSV with a header line.
	FormatCSV
)

// The --format flag and its values.
const (
	// FormatParamStr is the name of the flag which selects the format.
	FormatParamStr = "format"
	// FormatTableStr selects FormatTable.
	FormatTableStr = "table"
	// FormatJSONStr selects FormatJSON.
	FormatJSONStr = "json"
	// FormatCSVStr selects FormatCSV.
	FormatCSVStr = "csv"
)

func allSupportedFormats() []string {
	return []string{
		FormatTableStr,
		FormatJSONStr,
		FormatCSVStr,
	}
}

var formatMap = map[string]Format{
	FormatTableStr: FormatTable,
	FormatJSONStr:  FormatJSON,
	FormatCSVStr:   FormatCSV,
}

// NewFormatFromString returns the format from the string.
func NewFormatFromString(s string) (Format, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if format, ok := formatMap[s]; ok {
		return format, nil
	}
	return FormatTable, fmt.Errorf("invalid format: %s", s)
}

// String returns the string representation of the format.
func (f Format) String() string {
	for k, v := range formatMap {
		if v == f {
			return k
		}
	}
	return "unknown"
}
