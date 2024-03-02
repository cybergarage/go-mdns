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

// Record represents a record interface.
type Record interface {
	// SetName sets the resource record name.
	SetName(name string) Record
	// Name returns the resource record name.
	Name() string
	// SetType sets the resource record type.
	SetType(t Type) Record
	// Type returns the resource record type.
	Type() Type
	// SetClass sets the resource record class.
	SetClass(c Class) Record
	// Class returns the resource record class.
	Class() Class
	// SetTTL sets the TTL second.
	SetTTL(ttl uint) Record
	// TTL returns the TTL second.
	TTL() uint
	// SetUnicastResponse sets the unicast response flag.
	SetUnicastResponse(flag bool) Record
	// UnicastResponse returns the unicast response flag.
	UnicastResponse() bool
	// SetData sets the  record data.
	SetData(data []byte) Record
	// Data returns the record data.
	Data() []byte
	// Content returns a string representation to the record data.
	Content() string
	// Bytes returns the binary representation.
	Bytes() []byte
}
