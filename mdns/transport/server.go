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

package transport

import (
	"net"
)

// A Server represents a server.
type Server struct {
	Interface *net.Interface
}

// NewServer returns a new UnicastServer.
func NewServer() *Server {
	server := &Server{
		Interface: nil,
	}
	return server
}

// SetBoundInterface sets a bind interface.
func (server *Server) SetBoundInterface(i *net.Interface) {
	server.Interface = i
}

// GetBoundInterface return a bind interface.
func (server *Server) GetBoundInterface() *net.Interface {
	return server.Interface
}

// GetBoundAddresses returns the listen addresses.
func (server *Server) GetBoundAddresses() []string {
	boundAddrs := make([]string, 0)
	ifAddr, err := GetInterfaceAddress(server.Interface)
	if err != nil {
		return boundAddrs
	}
	boundAddrs = append(boundAddrs, ifAddr)
	return boundAddrs
}
