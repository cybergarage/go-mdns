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

// RFC 6763: 6. Data Syntax for DNS-SD TXT Records
// https://www.rfc-editor.org/rfc/rfc6763#section-6

// Attribute represents a DNS-SD TXT record key/value pair.
type Attribute interface {
	// Name returns the attribute name (key).
	Name() string
	// Value returns the attribute value. The value is empty when the
	// attribute has no value.
	Value() string
	// HasValue returns true if the attribute string has the '=' separator,
	// otherwise false. RFC 6763 (6.1) distinguishes an attribute with an
	// empty value ("key=") from a boolean attribute which is only present
	// ("key").
	HasValue() bool
	// String returns the attribute string.
	String() string
}
