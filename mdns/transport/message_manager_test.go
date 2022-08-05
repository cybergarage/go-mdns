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

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type testMessageManager struct {
	*MessageManager
	FromPort                int
	FromPacketType          int
	lastNotificationMessage *protocol.Message
}

func newTestMessage(tid uint) (*protocol.Message, error) {
	tidBytes := make([]byte, 2)
	encoding.IntegerToByte(tid, tidBytes)

	testMessageBytes := []byte{
		protocol.EHD1Echonet,
		protocol.EHD2Format1,
		tidBytes[0], tidBytes[1],
		0xA0, 0xB0, 0xC0,
		0xD0, 0xE0, 0xF0,
		protocol.ESVWriteReadRequest,
		3,
		1, 1, 'a',
		2, 2, 'b', 'c',
		3, 3, 'c', 'd', 'e',
	}

	return protocol.NewMessageWithBytes(testMessageBytes)
}

func isTestMessage(msg *protocol.Message) bool {
	return msg.IsESV(protocol.ESVWriteReadRequest)
}

// NewMessageManager returns a new message manager.
func newTestMessageManager() *testMessageManager {
	mgr := &testMessageManager{
		MessageManager:          NewMessageManager(),
		FromPort:                0,
		FromPacketType:          protocol.UnknownPacket,
		lastNotificationMessage: nil,
	}
	return mgr
}

func (mgr *testMessageManager) ProtocolMessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	// log.Trace("ProtocolMessageReceived (R) : %s", msg.String())

	if isTestMessage(msg) {
		copyMsg, err := protocol.NewMessageWithMessage(msg)
		if err == nil {
			// log.Trace("ProtocolMessageReceived (U) : %s", copyMsg.String())
			mgr.lastNotificationMessage = copyMsg
		}
	}

	return nil, nil
}

func testMulticastMessagingWithRunningManagers(t *testing.T, mgrs []*testMessageManager) {
	t.Helper()

	srcMgrs := []*testMessageManager{mgrs[0], mgrs[1]}
	dstMgrs := []*testMessageManager{mgrs[1], mgrs[0]}

	for n := 0; n < len(srcMgrs); n++ {
		srcMgr := srcMgrs[n]
		srcMgr.FromPacketType = protocol.UnknownPacket

		dstMgr := dstMgrs[n]
		dstMgr.FromPacketType = protocol.MulticastPacket
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

		log.Trace("CMP(M) : %s ?= %s", msg.String(), dstLastMsg.String())

		if !msg.Equals(dstLastMsg) {
			log.Trace("CMP(M) : %s != %s", msg.String(), dstLastMsg.String())
			t.Errorf("CMP(M) : %s != %s", msg.String(), dstLastMsg.String())
			continue
		}
	}
}
