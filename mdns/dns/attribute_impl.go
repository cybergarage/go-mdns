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

import (
	"fmt"
	"strings"
)

// RFC 6763: 6. Data Syntax for DNS-SD TXT Records
// https://www.rfc-editor.org/rfc/rfc6763#section-6

// attrImpl represents a DNS attribute.
type attrImpl struct {
	name     string
	value    string
	hasValue bool
}

// NewAttribute returns a new attribute instance.
func NewAttribute() Attribute {
	return newAttribute()
}

func newAttribute() *attrImpl {
	return &attrImpl{
		name:     "",
		value:    "",
		hasValue: false,
	}
}

// NewAttributeFromString returns a new attribute instance from the specified string.
func NewAttributeFromString(str string) (Attribute, error) {
	attr := newAttribute()
	return attr, attr.parse(str)
}

// parse parses the attribute string.
//
// RFC 6763: 6.4. Rules for Keys in DNS-SD Key/Value Pairs
// The key MUST be at least one character. The characters of a key MUST be
// printable US-ASCII values (0x20-0x7E), excluding '=' (0x3D).
//
// RFC 6763: 6.3. Rules for Values in DNS-SD Key/Value Pairs
// If there is no '=' in a DNS-SD TXT record string, then it is a boolean
// attribute, simply identified as being present, with no value. If the '=' is
// present, then everything after the first '=' is the value, so the value may
// contain '=' characters.
func (attr *attrImpl) parse(str string) error {
	if len(str) == 0 {
		return fmt.Errorf("attribute (%s) is %w", str, ErrInvalid)
	}

	name, value, hasValue := strings.Cut(str, "=")
	if len(name) == 0 {
		return fmt.Errorf("attribute (%s) is %w", str, ErrInvalid)
	}

	attr.name = name
	attr.value = value
	attr.hasValue = hasValue

	return nil
}

// Name returns the attribute name.
func (attr *attrImpl) Name() string {
	return attr.name
}

// Value returns the attribute value.
func (attr *attrImpl) Value() string {
	return attr.value
}

// HasValue returns true if the attribute has a value, otherwise false.
func (attr *attrImpl) HasValue() bool {
	return attr.hasValue
}

// String returns the attribute string.
func (attr *attrImpl) String() string {
	if !attr.hasValue {
		return attr.name
	}
	return attr.name + "=" + attr.value
}
