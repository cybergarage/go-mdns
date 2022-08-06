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
	"encoding/hex"
	"io"
)

// Message represents a protocol message.
type Message struct {
	*Header
}

// NewRequestMessage returns a request message instance.
func NewRequestMessage() *Message {
	msg := &Message{
		Header: NewRequestHeader(),
	}
	return msg
}

// NewMessageWithReader returns a message instance with the specified reader.
func NewMessageWithReader(reader io.Reader) (*Message, error) {
	msg := NewRequestMessage()
	if err := msg.Parse(reader); err != nil {
		return nil, err
	}
	return msg, nil
}

// NewMessageWithBytes returns a message instance with the specified bytes.
func NewMessageWithBytes(msgBytes []byte) (*Message, error) {
	return NewMessageWithReader(bytes.NewReader(msgBytes))
}

// Parse parses the specified reader.
func (msg *Message) Parse(reader io.Reader) error {
	if err := msg.Header.Parse(reader); err != nil {
		return err
	}
	return nil
}

// Equals returns true if the message is same as the specified message, otherwise false.
func (msg *Message) Equals(other *Message) bool {
	return msg.Header.Equals(other.Header)
}

// Copy returns the copy message instance.
func (msg *Message) Copy() *Message {
	return &Message{
		Header: NewHeaderWithBytes(msg.Header.bytes),
	}
}

// Bytes returns the binary representation.
func (msg *Message) Bytes() []byte {
	return msg.Header.Bytes()
}

// String returns the string representation.
func (msg *Message) String() string {
	return hex.EncodeToString(msg.Header.bytes)
}
