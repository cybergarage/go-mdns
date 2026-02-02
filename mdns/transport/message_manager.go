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

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A MessageManager represents a multicast server list.
type MessageManager struct {
	*MulticastManager
	*UnicastManager
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{
		MulticastManager: NewMulticastManager(),
		UnicastManager:   NewUnicastManager(),
	}
	return mgr
}

// SetMessageProcessor sets the message processor.
func (mgr *MessageManager) SetMessageProcessor(processor dns.MessageProcessor) {
	mgr.MulticastManager.SetMessageProcessor(processor)
	mgr.UnicastManager.SetMessageProcessor(processor)
}

// AnnounceMessage sends a message to the multicast address.
func (mgr *MessageManager) AnnounceMessage(msg dns.Message) error {
	return mgr.UnicastManager.AnnounceMessage(msg)
}

// Start starts this server.
func (mgr *MessageManager) Start() error {
	starter := []func() error{
		mgr.UnicastManager.Start,
		mgr.MulticastManager.Start,
	}
	var err error
	for _, startFunc := range starter {
		if startErr := startFunc(); startErr != nil {
			err = errors.Join(err, startErr)
		}
	}
	return err
}

// Stop stops this server.
func (mgr *MessageManager) Stop() error {
	stopper := []func() error{
		mgr.MulticastManager.Stop,
		mgr.UnicastManager.Stop,
	}
	var err error
	for _, stopFunc := range stopper {
		if stopErr := stopFunc(); stopErr != nil {
			err = errors.Join(err, stopErr)
		}
	}
	return err
}
