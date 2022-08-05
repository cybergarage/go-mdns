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
	"testing"
	"time"

	"github.com/cybergarage/go-mdns/mdns/protocol"
)

type testMessageManager struct {
	*MessageManager
	lastNotificationMessage *protocol.Message
}

func newTestMessage(tid uint) (*protocol.Message, error) {
	testMessageBytes := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01}
	return protocol.NewMessageWithBytes(testMessageBytes)
}

func isTestMessage(msg *protocol.Message) bool {
	return true
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

func testMulticastMessagingWithRunningManagers(t *testing.T, mgrs []*testMessageManager) {
	t.Helper()

	srcMgrs := []*testMessageManager{mgrs[0], mgrs[1]}
	dstMgrs := []*testMessageManager{mgrs[1], mgrs[0]}

	for n := 0; n < len(srcMgrs); n++ {
		dstMgr := dstMgrs[n]
		dstMgr.lastNotificationMessage = nil

		msg, err := newTestMessage(uint(n | 0xF0))
		if err != nil {
			t.Error(err)
			continue
		}

		// err = srcMgr.AnnounceMessage(msg)
		// if err != nil {
		// 	t.Error(err)
		// 	continue
		// }

		time.Sleep(time.Second)

		dstLastMsg := dstMgr.lastNotificationMessage
		if dstLastMsg == nil {
			t.Errorf("%s !=", msg.String())
			continue
		}

		if !msg.Equals(dstLastMsg) {
			t.Errorf("CMP(M) : %s != %s", msg.String(), dstLastMsg.String())
			continue
		}
	}
}
