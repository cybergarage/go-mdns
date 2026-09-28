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
	"errors"
	"io"
	"net"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A MulticastServer represents a multicast server.
type MulticastServer struct {
	*Server
	*MulticastSocket
	channel   chan any
	processor dns.MessageProcessor
}

// NewMulticastServer returns a new MulticastServer.
func NewMulticastServer() *MulticastServer {
	server := &MulticastServer{
		Server:          NewServer(),
		MulticastSocket: NewMulticastSocket(),
		channel:         nil,
		processor:       nil,
	}
	return server
}

// SetMessageProcessor sets the message processor.
func (server *MulticastServer) SetMessageProcessor(processor dns.MessageProcessor) {
	server.processor = processor
}

// Start starts this server.
func (server *MulticastServer) Start(ifi *net.Interface, ifaddr string) error {
	if err := server.MulticastSocket.Bind(ifi, ifaddr); err != nil {
		return err
	}
	server.channel = make(chan any)
	go handleMulticastConnection(server, server.channel)
	return nil
}

// Stop stops this server.
func (server *MulticastServer) Stop() error {
	if err := server.MulticastSocket.Close(); err != nil {
		return err
	}
	server.SetListenInterface(nil)
	return nil
}

func handleMulticastRequestMessage(server *MulticastServer, reqMsg dns.Message) {
	if server.processor == nil {
		return
	}
	resMsg, err := server.processor(reqMsg)
	if err != nil || resMsg == nil {
		return
	}
	if addr, port, ok := unicastResponseAddr(reqMsg); ok {
		if _, err := server.SendMessage(addr, port, resMsg); err != nil {
			log.Debugf("Failed to send unicast response: %s", err)
		}
		return
	}
	if err := server.AnnounceMessage(resMsg); err != nil {
		log.Debugf("Failed to send multicast response: %s", err)
	}
}

// unicastResponseAddr returns where the response to reqMsg is sent directly
// rather than multicast: to the source of a legacy unicast query, sent from
// a port other than 5353 (RFC 6762, 6.7), and to the source of a query
// whose questions all ask for a unicast response (5.4).
func unicastResponseAddr(reqMsg dns.Message) (string, int, bool) {
	from := reqMsg.From()
	if from == nil || from.IP() == nil {
		return "", 0, false
	}
	if from.Port() == Port && !allQuestionsUnicast(reqMsg) {
		return "", 0, false
	}
	addr := from.IP().String()
	if zone := from.Zone(); zone != "" {
		addr += "%" + zone
	}
	return addr, from.Port(), true
}

func allQuestionsUnicast(msg dns.Message) bool {
	questions := msg.Questions()
	if len(questions) == 0 {
		return false
	}
	for _, q := range questions {
		if !q.UnicastResponse() {
			return false
		}
	}
	return true
}

func handleMulticastConnection(server *MulticastServer, cancel chan any) {
	defer server.Socket.Close()
	for {
		select {
		case <-cancel:
			return
		default:
			msg, err := server.MulticastSocket.ReadMessage()
			if err != nil {
				if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
					return
				}
				log.Debugf("Failed to read multicast message: %s", err)
				continue
			}
			go handleMulticastRequestMessage(server, msg)
		}
	}
}
