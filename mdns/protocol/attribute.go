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

// RFC1464: Using the Domain Name System To Store Arbitrary String Attributes
// https://datatracker.ietf.org/doc/html/rfc1464

// Attribute represents a DNS attribute.
type Attribute struct {
	name  string
	value string
}

// NewAttribute returns a new attribute innstance.
func NewAttribute() *Attribute {
	return &Attribute{
		name:  "",
		value: "",
	}
}

func NewAttributeWithString(str string) (*Attribute, error) {
	attr := NewAttribute()
	err := attr.parse(str)
	return attr, err
}

// Parse parses the attribute string.
func (attr *Attribute) parse(str string) error {
	// TODO : Implement this
	return nil
}

// Name returns the attribute name.
func (attr *Attribute) Name() string {
	return attr.name
}

// Value returns the attribute value.
func (attr *Attribute) Value() string {
	return attr.value
}

// String returns the attribute string.
func (attr *Attribute) String() string {
	return attr.name + "=" + attr.value
}
