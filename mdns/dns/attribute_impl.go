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

// RFC1464: Using the Domain Name System To Store Arbitrary String Attributes
// https://datatracker.ietf.org/doc/html/rfc1464

// attrImpl represents a DNS attribute.
type attrImpl struct {
	name  string
	value string
}

// NewAttribute returns a new attribute instance.
func NewAttribute() Attribute {
	return &attrImpl{
		name:  "",
		value: "",
	}
}

func newAttribute() *attrImpl {
	return &attrImpl{
		name:  "",
		value: "",
	}
}

// NewAttributeFromString returns a new attribute instance from the specified string.
func NewAttributeFromString(str string) (Attribute, error) {
	attr := newAttribute()
	return attr, attr.parse(str)
}

// Parse parses the attribute string.
func (attr *attrImpl) parse(str string) error {
	vars := strings.Split(str, "=")
	if len(vars) != 2 {
		return fmt.Errorf("attribute (%s) is %w", str, ErrInvalid)
	}
	attr.name = vars[0]
	attr.value = vars[1]
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

// String returns the attribute string.
func (attr *attrImpl) String() string {
	return attr.name + "=" + attr.value
}
