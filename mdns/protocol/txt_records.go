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

// TXTRecords is a list of TXTRecord.
type TXTRecords []*TXTRecord

// LookupAttribute returns the attribute with the specified name.
func (txts TXTRecords) LookupAttribute(name string) (*Attribute, bool) {
	for _, txt := range txts {
		attr, ok := txt.LookupAttribute(name)
		if ok {
			return attr, true
		}
	}
	return nil, false
}
