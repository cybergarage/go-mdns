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
	"github.com/cybergarage/go-mdns/mdns/dns"
)

// A MessageManager represents a multicast server list.
type MessageManager struct {
	msgProcessor dns.MessageProcessor
	*MulticastManager
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{
		msgProcessor:     nil,
		MulticastManager: NewMulticastManager(),
	}
	return mgr
}

// SetMessageProcessor sets the message processor.
func (mgr *MessageManager) SetMessageProcessor(h dns.MessageProcessor) {
	mgr.SetHandler(h)
	mgr.msgProcessor = h
}

// MessageProcessor returns the message processor.
func (mgr *MessageManager) MessageProcessor() dns.MessageProcessor {
	return mgr.msgProcessor
}
