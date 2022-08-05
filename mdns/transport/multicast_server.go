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

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A MulticastServer represents a multicast server.
type MulticastServer struct {
	*Server
	Socket        *MulticastSocket
	Channel       chan interface{}
	Handler       MulticastHandler
	UnicastServer *UnicastServer
}

// NewMulticastServer returns a new MulticastServer.
func NewMulticastServer() *MulticastServer {
	server := &MulticastServer{
		Server:        NewServer(),
		Socket:        NewMulticastSocket(),
		Channel:       nil,
		Handler:       nil,
		UnicastServer: nil,
	}
	return server
}

// SetHandler set a listener.
func (server *MulticastServer) SetHandler(l MulticastHandler) {
	server.Handler = l
}

// SetUnicastServer set a unicast server to response received messages.
func (server *MulticastServer) SetUnicastServer(s *UnicastServer) {
	server.UnicastServer = s
}

// Start starts this server.
func (server *MulticastServer) Start(ifi *net.Interface) error {
	if err := server.Socket.Bind(ifi); err != nil {
		return err
	}
	server.SetBoundInterface(ifi)
	server.Channel = make(chan interface{})
	go handleMulticastConnection(server, server.Channel)
	return nil
}

// Stop stops this server.
func (server *MulticastServer) Stop() error {
	if err := server.Socket.Close(); err != nil {
		return err
	}
	server.SetBoundInterface(nil)
	return nil
}

func handleMulticastRequestMessage(server *MulticastServer, reqMsg *protocol.Message) {
	server.Socket.outputReadLog(log.LevelTrace, logSocketTypeUDPMulticast, reqMsg.From.String(), reqMsg.String(), reqMsg.Size())

	if server.Handler == nil {
		return
	}

	resMsg, err := server.Handler.ProtocolMessageReceived(reqMsg)
	if server.UnicastServer == nil || err != nil || resMsg == nil {
		return
	}

	server.UnicastServer.UDPSocket.ResponseMessageForRequestMessage(reqMsg, resMsg)
}

func handleMulticastConnection(server *MulticastServer, cancel chan interface{}) {
	defer server.Socket.Close()
	for {
		select {
		case <-cancel:
			return
		default:
			reqMsg, err := server.Socket.ReadMessage()
			if err != nil {
				break
			}
			reqMsg.SetPacketType(protocol.MulticastPacket)

			go handleMulticastRequestMessage(server, reqMsg)
		}
	}
}
