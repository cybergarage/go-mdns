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
	"testing"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

type testMulticastServer struct {
	*MulticastServer
	lastMessage *dns.Message
}

// NewMessageManager returns a new message manager.
func newTestMulticastServer() *testMulticastServer {
	server := &testMulticastServer{
		MulticastServer: NewMulticastServer(),
		lastMessage:     nil,
	}
	server.SetHandler(server)
	return server
}

func (server *testMulticastServer) MessageReceived(msg *dns.Message) (*dns.Message, error) {
	if isTestMessage(msg) {
		server.lastMessage = msg.Copy()
	}
	return nil, nil
}

func testMulticastServerWithInterface(t *testing.T, ifi *net.Interface, ifaddr string) {
	t.Helper()

	server := newTestMulticastServer()

	// Start server

	err := server.Start(ifi, ifaddr)
	if err != nil {
		t.Error(err)
		return
	}

	// Stop server

	defer func() {
		err := server.Stop()
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(time.Second)

	// Send a test message

	now := time.Now()
	msg, err := newTestMessage(uint(now.Unix()))
	if err != nil {
		t.Error(err)
		return
	}

	toAddr := MulticastIPv4Address
	if IsIPv6Address(ifaddr) {
		toAddr = MulticastIPv6Address
	}
	nSent, err := server.SendMessage(toAddr, Port, msg)
	if err != nil {
		t.Error(err)
		return
	}

	if msgBytes := msg.Bytes(); nSent != len(msgBytes) {
		t.Errorf("%d != %d", nSent, len(msgBytes))
		return
	}

	// Wait a test message

	// time.Sleep(time.Second)

	// if !msg.Equal(server.lastMessage) {
	// t.Errorf("%s != %s", msg.String(), server.lastMessage.String())
	// }
}

func TestMulticastServerWithInterface(t *testing.T) {
	log.EnableStdoutDebug(true)
	defer log.EnableStdoutDebug(false)

	ifis, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
		return
	}

	for _, ifi := range ifis {
		t.Run(ifi.Name, func(t *testing.T) {
			ifaddrs, err := GetInterfaceAddresses(ifi)
			if err != nil {
				t.Error(err)
				return
			}
			for _, ifaddr := range ifaddrs {
				t.Run(ifaddr, func(t *testing.T) {
					testMulticastServerWithInterface(t, ifi, ifaddr)
				})
			}
		})
	}
}
