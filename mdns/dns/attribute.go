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

// RFC1464: Using the Domain Name System To Store Arbitrary String Attributes
// https://datatracker.ietf.org/doc/html/rfc1464

// Attribute represents a DNS attribute.
type Attribute interface {
	// Name returns the attribute name.
	Name() string
	// Value returns the attribute value.
	Value() string
	// String returns the attribute string.
	String() string
}
