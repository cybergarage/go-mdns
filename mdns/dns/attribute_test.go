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
	"testing"
)

// RFC 6763: 6.3. Rules for Values in DNS-SD Key/Value Pairs.
func TestAttributeFromString(t *testing.T) {
	tests := []struct {
		str      string
		name     string
		value    string
		hasValue bool
		isValid  bool
	}{
		{"SII=5000", "SII", "5000", true, true},
		{"VP=65521+32769", "VP", "65521+32769", true, true},
		{"CM=2", "CM", "2", true, true},
		// Everything after the first '=' is the value.
		{"key=a=b", "key", "a=b", true, true},
		// An attribute with an empty value.
		{"key=", "key", "", true, true},
		// A boolean attribute which is only present.
		{"key", "key", "", false, true},
		// An attribute without a key is invalid.
		{"=value", "", "", false, false},
		{"", "", "", false, false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			attr, err := NewAttributeFromString(test.str)
			if !test.isValid {
				if err == nil {
					t.Errorf("%s : expected an error", test.str)
				}
				return
			}
			if err != nil {
				t.Error(err)
				return
			}
			if attr.Name() != test.name {
				t.Errorf("name %s != %s", attr.Name(), test.name)
			}
			if attr.Value() != test.value {
				t.Errorf("value %s != %s", attr.Value(), test.value)
			}
			if attr.HasValue() != test.hasValue {
				t.Errorf("hasValue %t != %t", attr.HasValue(), test.hasValue)
			}
			if attr.String() != test.str {
				t.Errorf("string %s != %s", attr.String(), test.str)
			}
		})
	}
}

// RFC 6763: 6.4. Rules for Keys in DNS-SD Key/Value Pairs.
func TestAttributesLookup(t *testing.T) {
	attrs, err := NewAttributesFromStrings([]string{"SII=5000", "sii=9999", "D=3840", "CM"})
	if err != nil {
		t.Fatal(err)
	}

	// The keys are case insensitive.
	for _, name := range []string{"SII", "sii", "Sii"} {
		attr, ok := attrs.LookupAttribute(name)
		if !ok {
			t.Errorf("%s is not found", name)
			continue
		}
		// The first occurrence is used, and the later ones are ignored.
		if attr.Value() != "5000" {
			t.Errorf("%s = %s != 5000", name, attr.Value())
		}
	}

	attr, ok := attrs.LookupAttribute("CM")
	if !ok {
		t.Fatal("CM is not found")
	}
	if attr.HasValue() {
		t.Errorf("CM should not have a value")
	}

	if attrs.HasAttribute("PH") {
		t.Errorf("PH should not be found")
	}
}
