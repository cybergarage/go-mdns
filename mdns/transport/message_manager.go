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
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A MessageManager represents a multicast server list.
type MessageManager struct {
	messageHandler protocol.MessageHandler
	multicastMgr   *MulticastManager
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{
		messageHandler: nil,
		multicastMgr:   NewMulticastManager(),
	}
	return mgr
}

// SetMessageHandler set a listener to all managers.
func (mgr *MessageManager) SetMessageHandler(h protocol.MessageHandler) {
	mgr.multicastMgr.SetHandler(h)
	mgr.messageHandler = h
}

// GetMessageHandler returns the listener of the manager.
func (mgr *MessageManager) GetMessageHandler() protocol.MessageHandler {
	return mgr.messageHandler
}

// Start starts all transport managers.
func (mgr *MessageManager) Start() error {
	err := mgr.Stop()
	if err != nil {
		return err
	}
	err = mgr.multicastMgr.Start()
	if err != nil {
		mgr.Stop()
		return err
	}
	return nil
}

// Stop stops all transport managers.
func (mgr *MessageManager) Stop() error {
	return mgr.multicastMgr.Stop()
}

// IsRunning returns true whether the local managers are running, otherwise false.
func (mgr *MessageManager) IsRunning() bool {
	return mgr.multicastMgr.IsRunning()
}
