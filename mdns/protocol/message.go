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
	Questions
}

// NewMessage returns a nil message instance.
func NewMessage() *Message {
	msg := &Message{
		Header:    NewHeader(),
		Questions: Questions{},
	}
	return msg
}

// NewRequestMessage returns a request message instance.
func NewRequestMessage() *Message {
	msg := &Message{
		Header:    NewRequestHeader(),
		Questions: Questions{},
	}
	return msg
}

// NewResponseMessage returns a response message instance.
func NewResponseMessage() *Message {
	msg := &Message{
		Header:    NewResponseHeader(),
		Questions: Questions{},
	}
	return msg
}

// NewMessageWithReader returns a message instance with the specified reader.
func NewMessageWithReader(reader io.Reader) (*Message, error) {
	msg := NewMessage()
	if err := msg.Parse(reader); err != nil {
		return nil, err
	}
	return msg, nil
}

// NewMessageWithBytes returns a message instance with the specified bytes.
func NewMessageWithBytes(msgBytes []byte) (*Message, error) {
	return NewMessageWithReader(bytes.NewReader(msgBytes))
}

// AddQuestion adds the specified question into the message.
func (msg *Message) AddQuestion(q *Question) {
	msg.Questions = append(msg.Questions, q)
}

// Parse parses the specified reader.
func (msg *Message) Parse(reader io.Reader) error {
	if err := msg.Header.Parse(reader); err != nil {
		return err
	}
	// Parses questions
	for n := 0; n < int(msg.QD()); n++ {
		q, err := NewQuestionWithReader(reader)
		if err != nil {
			return nil
		}
		msg.AddQuestion(q)
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
		Header:    NewHeaderWithBytes(msg.Header.bytes),
		Questions: msg.Questions,
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
