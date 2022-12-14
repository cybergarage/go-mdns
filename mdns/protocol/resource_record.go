// Copyright (C) 2022 Satoshi Konno All rights reserved.
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

// ResourceRecord represents a resource record interface.
type ResourceRecord interface {
	// Name returns the resource record name.
	Name() string
	// Type returns the resource record type.
	Type() Type
	// Class returns the resource record class.
	Class() Class
	// UnicastResponse returns the unicast response flag.
	UnicastResponse() bool
	// TTL returns the TTL second.
	TTL() uint
	// Data returns the resource record data.
	Data() []byte
	// RequestBytes returns only the binary representation of the request fields.
	RequestBytes() []byte
	// ResponseBytes returns only the binary representation of the all fields.
	ResponseBytes() []byte
	// Bytes returns the binary representation.
	Bytes() []byte
	// Equal returns true if this record is equal to  the specified resource record. otherwise false.
	Equal(res ResourceRecord) bool
}
