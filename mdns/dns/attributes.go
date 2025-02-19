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
	"strings"
)

type Attributes []*Attribute

// NewAttributes returns a new attributes instance.
func NewAttributes() Attributes {
	return Attributes{}
}

// NewAttributesFromStrings returns a new attributes instance from the specified strings.
func NewAttributesFromStrings(strs []string) (Attributes, error) {
	attrs := NewAttributes()
	for _, str := range strs {
		attr, err := NewAttributeFromString(str)
		if err != nil {
			continue
		}
		attrs = append(attrs, attr)
	}
	return attrs, nil
}

// LookupAttribute returns the attribute with the specified name.
func (attrs Attributes) LookupAttribute(name string) (*Attribute, bool) {
	for _, attr := range attrs {
		if attr.Name() == name {
			return attr, true
		}
	}
	return nil, false
}

// HasAttribute returns true if this instance has the specified attribute.
func (attrs Attributes) HasAttribute(name string) bool {
	_, ok := attrs.LookupAttribute(name)
	return ok
}

// String returns the attribute string.
func (attrs Attributes) String() string {
	strs := []string{}
	for _, attr := range attrs {
		strs = append(strs, attr.String())
	}
	return strings.Join(strs, " ")
}
