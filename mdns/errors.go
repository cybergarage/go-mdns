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

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalid indicates that a specified value is invalid.
	ErrInvalid = errors.New("invalid")
	// ErrNotFound indicates that a queried name is not found.
	ErrNotFound = errors.New("not found")
	// ErrNotRunning indicates that the server is not running: a service is
	// registered before Server.Start, or the server stops while it probes.
	ErrNotRunning = errors.New("mdns: server is not running")
	// ErrConflict indicates that another node on the link holds a name of
	// a service (RFC 6762, 8.1 and 9). A *ConflictError wraps it.
	ErrConflict = errors.New("mdns: name conflict")
	// ErrDeregistered indicates that a service was deregistered, or
	// replaced by a later registration of the same instance, while it was
	// being probed.
	ErrDeregistered = errors.New("mdns: service deregistered")
)

// ConflictError reports that another node on the link holds one of the names
// of a service, its instance name or its host name, so the server does not
// publish the service.
//
// The server does not rename a service: the application decides on a new
// name, such as a new random instance name for a Matter commissionable node,
// and registers the service again.
type ConflictError struct {
	// Service is the service whose name conflicts.
	Service *LocalService
	// Name is the conflicting name with its domain: the instance name, such
	// as "665F6E75B5D3A9C2._matterc._udp.local", or the host name, such as
	// "B75AFB458ECD6D6F.local".
	Name string
}

// Error returns the error message.
func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s: %s", ErrConflict, e.Name)
}

// Unwrap returns ErrConflict, so that errors.Is(err, ErrConflict) holds.
func (e *ConflictError) Unwrap() error {
	return ErrConflict
}

// IsInstanceConflict reports whether the instance name of the service
// conflicts.
func (e *ConflictError) IsInstanceConflict() bool {
	return e.Service != nil && strings.EqualFold(e.Name, e.Service.FullName())
}

// IsHostConflict reports whether the host name of the service conflicts.
func (e *ConflictError) IsHostConflict() bool {
	return e.Service != nil && strings.EqualFold(e.Name, e.Service.HostName())
}
