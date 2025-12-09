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

// SRVRecord represents a SRV record.
// RFC 2782: A DNS RR for specifying the location of services (DNS SRV).
// https://www.rfc-editor.org/rfc/rfc2782
type SRVRecord interface {
	Record
	// Service returns the service name.
	Service() string
	// Proto returns the protocol name.
	Proto() string
	// Name returns the domain name.
	Name() string
	// Priority returns the priority of the target host.
	Priority() uint
	// Weight returns a relative weight for records with the same priority.
	Weight() uint
	// Port returns the port on this target host of this service.
	Port() uint
	// Target returns the canonical hostname of the machine providing the service.
	Target() string
	// Content returns a string representation to the record data.
	Content() string
}
