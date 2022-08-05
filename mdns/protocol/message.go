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

package protocol

import (
	"bytes"
	"io"
)

// Message represents a protocol message.
type Message struct {
	*Header
}

// NewMessage returns a message instance.
func NewMessage() *Message {
	msg := &Message{
		Header: NewHeader(),
	}
	return msg
}

// NewMessage returns a message instance.
func NewMessageWithBytes(msgBytes []byte) (*Message, error) {
	msg := NewMessage()
	if err := msg.Parse(bytes.NewReader(msgBytes)); err != nil {
		return nil, err
	}
	return msg, nil
}

// Parse parses the specified reader.
func (msg *Message) Parse(reader io.Reader) error {
	if err := msg.Header.Parse(reader); err != nil {
		return err
	}
	return nil
}

// Copy returns the copy message instance.
func (msg *Message) Copy() *Message {
	return &Message{
		Header: NewHeaderWithBytes(msg.Header.bytes),
	}
}
