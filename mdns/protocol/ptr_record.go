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

// PTRRecord represents a PTR record.
type PTRRecord struct {
	*resourceRecord
}

// newPTRRecordWithResourceRecord returns a new PTR record innstance.
func newPTRRecordWithResourceRecord(res *resourceRecord) *PTRRecord {
	return &PTRRecord{
		resourceRecord: res,
	}
}

// DomainName returns the resource domain name.
func (ptr *PTRRecord) DomainName() string {
	if len(ptr.data) < 1 {
		return ""
	}
	return string(ptr.data[1:])
}