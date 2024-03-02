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
	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/cybergarage/go-mdns/mdns/transport"
)

// Server represents a server node instance.
type Server struct {
	*transport.MessageManager
	*services
	userListener MessageListener
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		MessageManager: transport.NewMessageManager(),
		services:       newServices(),
		userListener:   nil,
	}
	server.SetMessageHandler(server)
	return server
}

// Set sets a message listner to listen raw protocol messages.
func (server *Server) SetListener(l MessageListener) {
	server.userListener = l
}

// Start starts the server instance.
func (server *Server) Start() error {
	if err := server.Stop(); err != nil {
		return err
	}
	return server.MessageManager.Start()
}

// Stop stops the server instance.
func (server *Server) Stop() error {
	return server.MessageManager.Stop()
}

// Restart restarts the server instance.
func (server *Server) Restart() error {
	if err := server.Stop(); err != nil {
		return err
	}
	return server.Start()
}

func (server *Server) MessageReceived(msg *dns.Message) (*dns.Message, error) {
	if server.userListener != nil {
		server.userListener.MessageReceived(msg)
	}

	if msg.IsResponse() {
		return nil, nil
	}

	return nil, nil
}
