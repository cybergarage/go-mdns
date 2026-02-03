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

package dns

// IsQueryWithUnicastResponse returns true if the message is a query with the unicast response bit set, otherwise false.
func (msg *message) IsQueryWithUnicastResponse() bool {
	if msg == nil {
		return false
	}
	if !msg.IsQuery() {
		return false
	}
	for _, q := range msg.Questions() {
		if q.IsUnicastResponse() {
			return true
		}
	}
	return false
}

// IsQueryAnswer returns true if the message is a response to a query, otherwise false.
func (msg *message) IsQueryAnswer(resMsg Message) bool {
	if msg == nil {
		return true
	}
	if resMsg == nil {
		return false
	}
	if !msg.IsQuery() || !resMsg.IsResponse() {
		return false
	}
	if msg.ID() != 0 && resMsg.ID() != 0 && msg.ID() != resMsg.ID() {
		return false
	}
	for _, q := range msg.Questions() {
		for _, rr := range resMsg.ResourceRecordSet() {
			if !rr.IsName(q.Name()) {
				continue
			}
			if !q.Class().Equal(rr.Class()) {
				continue
			}
			if q.Type() != ANY && q.Type() != rr.Type() {
				continue
			}
			return true
		}
	}
	return false
}
