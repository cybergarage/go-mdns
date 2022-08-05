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

	"github.com/cybergarage/go-mdns/mdns/protocol"
)

// A MulticastServer represents a multicast server.
type MulticastServer struct {
	*Server
	*MulticastSocket
	Channel chan interface{}
	Handler MulticastHandler
}

// NewMulticastServer returns a new MulticastServer.
func NewMulticastServer() *MulticastServer {
	server := &MulticastServer{
		Server:          NewServer(),
		MulticastSocket: NewMulticastSocket(),
		Channel:         nil,
		Handler:         nil,
	}
	return server
}

// SetHandler set a listener.
func (server *MulticastServer) SetHandler(handler MulticastHandler) {
	server.Handler = handler
}

// Start starts this server.
func (server *MulticastServer) Start(ifi *net.Interface) error {
	if err := server.MulticastSocket.Bind(ifi); err != nil {
		return err
	}
	server.SetBoundInterface(ifi)
	server.Channel = make(chan interface{})
	go handleMulticastConnection(server, server.Channel)
	return nil
}

// Stop stops this server.
func (server *MulticastServer) Stop() error {
	if err := server.MulticastSocket.Close(); err != nil {
		return err
	}
	server.SetBoundInterface(nil)
	return nil
}

func handleMulticastRequestMessage(server *MulticastServer, msg *protocol.Message) {
	if server.Handler == nil {
		return
	}
	server.Handler.MessageReceived(msg)
}

func handleMulticastConnection(server *MulticastServer, cancel chan interface{}) {
	defer server.Socket.Close()
	for {
		select {
		case <-cancel:
			return
		default:
			msg, err := server.MulticastSocket.ReadMessage()
			if err != nil {
				break
			}
			go handleMulticastRequestMessage(server, msg)
		}
	}
}
