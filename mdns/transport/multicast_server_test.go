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

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type testMulticastServer struct {
	*MulticastServer
	lastMessage *protocol.Message
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

func (server *testMulticastServer) ProtocolMessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	if isTestMessage(msg) {
		copyMsg, err := protocol.NewMessageWithMessage(msg)
		if err == nil {
			server.lastMessage = copyMsg
		}
	}

	return nil, nil
}

func testMulticastServerWithInterface(t *testing.T, ifi *net.Interface) {
	t.Helper()

	server := newTestMulticastServer()

	// Start server

	err := server.Start(ifi)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second)

	// Send a test message

	now := time.Now()
	msg, err := newTestMessage(uint(now.Unix()))
	if err != nil {
		t.Error(err)
		return
	}

	sock := NewUnicastUDPSocket()
	nSent, err := sock.SendMessage(MulticastAddress, Port, msg)
	if err != nil {
		t.Error(err)
	}

	if msgBytes := msg.Bytes(); nSent != len(msgBytes) {
		t.Errorf("%d != %d", nSent, len(msgBytes))
		return
	}

	// Wait a test message

	time.Sleep(time.Second)

	if !msg.Equals(server.lastMessage) {
		ifi, _ := server.MulticastServer.Socket.GetBoundInterface()
		t.Errorf("%v", ifi)
		t.Errorf("%s != %s", msg.String(), server.lastMessage.String())
	}

	// Stop server

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestMulticastServerWithInterface(t *testing.T) {
	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
		return
	}

	testMulticastServerWithInterface(t, ifs[0])
}

func TestMulticastServerWithNoInterface(t *testing.T) {
	testMulticastServerWithInterface(t, nil)
}
