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

package mdns

// ServiceEventType represents a change type of a discovered service.
type ServiceEventType int

const (
	// ServiceAdded is the event type for a service which is discovered for
	// the first time.
	ServiceAdded ServiceEventType = iota
	// ServiceUpdated is the event type for a service whose records have
	// changed, such as its addresses or its TXT attributes.
	ServiceUpdated
	// ServiceRemoved is the event type for a service which has gone away.
	// A service is removed when it announces a goodbye packet (RFC 6762,
	// 10.1), or when its records expire.
	ServiceRemoved
)

// String returns the string representation of the event type.
func (t ServiceEventType) String() string {
	switch t {
	case ServiceAdded:
		return "added"
	case ServiceUpdated:
		return "updated"
	case ServiceRemoved:
		return "removed"
	}
	return "unknown"
}

// ServiceEvent represents a change of a discovered service.
type ServiceEvent struct {
	// Type is the change type of the service.
	Type ServiceEventType
	// Service is the changed service.
	Service Service
}

// String returns the string representation of the event.
func (event ServiceEvent) String() string {
	return event.Type.String() + " " + event.Service.String()
}

// ServiceHandler is a handler function which is called when a browsed service
// is added, updated or removed.
type ServiceHandler func(ServiceEvent)
