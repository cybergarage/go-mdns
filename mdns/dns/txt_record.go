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

type TXTRecord interface {
	Record
	// Strings returns the resource text strings.
	Strings() []string
	// Attributes returns the resource attributes.
	Attributes() (Attributes, error)
	// LookupAttribute looks up the attribute by name.
	LookupAttribute(name string) (Attribute, bool)
	// Content returns a string representation to the record data.
	Content() string
}
