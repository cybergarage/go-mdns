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

// NSECRecord represents a NSEC record.
// RFC 4034: Resource Records for the DNS Security Extensions.
// https://www.rfc-editor.org/rfc/rfc4034
type NSECRecord interface {
	Record
	// Content returns a string representation to the record data.
	Content() string
}
