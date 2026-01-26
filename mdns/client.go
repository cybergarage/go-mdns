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
	"context"
)

// Client represents a client node instance.
type Client interface {
	// Start starts the client.
	Start() error
	// Stop stops the client.
	Stop() error
	// Restart restarts the client.
	Restart() error
	// RegisterHandler adds a message handler to the client.
	RegisterHandler(handler MessageHandler)
	// UnregisterHandler removes a message handler from the client.
	UnregisterHandler(handler MessageHandler)
	// Query sends a question message to the multicast address.
	Query(ctx context.Context, query Query) ([]Service, error)
}
