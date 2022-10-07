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
	"github.com/cybergarage/go-mdns/mdns/protocol"
)

type testMessageManager struct {
	*MessageManager
	lastNotificationMessage *protocol.Message
}

func newTestMessage(tid uint) (*protocol.Message, error) {
	msg := protocol.NewRequestMessage()
	return msg, nil
}

func isTestMessage(msg *protocol.Message) bool {
	return msg.ID() == 0x02
}

// newTestMessageManager returns a new message manager.
func newTestMessageManager() *testMessageManager {
	mgr := &testMessageManager{
		MessageManager:          NewMessageManager(),
		lastNotificationMessage: nil,
	}
	return mgr
}

func (mgr *testMessageManager) MessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	if isTestMessage(msg) {
		mgr.lastNotificationMessage = msg.Copy()
	}
	return nil, nil
}
