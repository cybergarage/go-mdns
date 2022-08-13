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
	"github.com/cybergarage/go-mdns/mdns/protocol"
	"github.com/cybergarage/go-mdns/mdns/transport"
)

// Client represents a client node instance.
type Client struct {
	*transport.MessageManager
	userListener MessageListener
}

// NewClient returns a new client instance.
func NewClient() *Client {
	client := &Client{
		MessageManager: transport.NewMessageManager(),
		userListener:   nil,
	}
	client.SetMessageHandler(client)
	return client
}

// Set sets a message listner to listen raw protocol messages.
func (client *Client) SetListener(l MessageListener) {
	client.userListener = l
}

// Start starts the client instance.
func (client *Client) Start() error {
	if err := client.Stop(); err != nil {
		return err
	}
	return client.MessageManager.Start()
}

// Stop stops the client instance.
func (client *Client) Stop() error {
	return client.MessageManager.Stop()
}

// Restart restarts the client instance.
func (client *Client) Restart() error {
	if err := client.Stop(); err != nil {
		return err
	}
	return client.Start()
}

// Query sends a question message to the multicast address.
func (client *Client) Query(q *Query) {
	msg := protocol.NewRequestMessage()
	msg.AddQuestion(protocol.NewQuestion().SetName(q.String()).SetType(protocol.PTR))
}

func (client *Client) MessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	if msg.IsQuery() {
		return nil, nil
	}
	return nil, nil
}
