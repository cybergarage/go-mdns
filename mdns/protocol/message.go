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
	Answers
	NameServers
	Additions
}

// NewMessage returns a nil message instance.
func NewMessage() *Message {
	msg := &Message{
		Header:      NewHeader(),
		Questions:   Questions{},
		Answers:     Answers{},
		NameServers: NameServers{},
		Additions:   Additions{},
	}
	return msg
}

// NewRequestMessage returns a request message instance.
func NewRequestMessage() *Message {
	msg := NewMessage()
	msg.Header = NewRequestHeader()
	return msg
}

// NewResponseMessage returns a response message instance.
func NewResponseMessage() *Message {
	msg := NewMessage()
	msg.Header = NewResponseHeader()
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
	msg.setQD(uint(len(msg.Questions)))
}

// AddAnswer adds the specified answer into the message.
func (msg *Message) AddAnswer(a *Answer) {
	msg.Answers = append(msg.Answers, a)
	msg.setAN(uint(len(msg.Answers)))
}

// AddNameServer adds the specified name server into the message.
func (msg *Message) AddNameServer(ns *NameServer) {
	msg.NameServers = append(msg.NameServers, ns)
	msg.setNS(uint(len(msg.NameServers)))
}

// AddAddition adds the specified additional record into the message.
func (msg *Message) AddAddition(a *Addition) {
	msg.Additions = append(msg.Additions, a)
	msg.setAR(uint(len(msg.Additions)))
}

// Parse parses the specified reader.
func (msg *Message) Parse(reader io.Reader) error {
	if err := msg.Header.Parse(reader); err != nil {
		return err
	}
	// Parses questions.
	for n := 0; n < int(msg.QD()); n++ {
		q, err := NewQuestionWithReader(reader)
		if err != nil {
			return nil
		}
		msg.Questions = append(msg.Questions, q)
	}
	// Parses answers.
	for n := 0; n < int(msg.AN()); n++ {
		a, err := NewResourceWithReader(reader)
		if err != nil {
			return nil
		}
		msg.Answers = append(msg.Answers, a)
	}
	// Parses name servers.
	for n := 0; n < int(msg.AN()); n++ {
		a, err := NewResourceWithReader(reader)
		if err != nil {
			return nil
		}
		msg.Answers = append(msg.Answers, a)
	}
	// Parses name servers.
	for n := 0; n < int(msg.NS()); n++ {
		ns, err := NewResourceWithReader(reader)
		if err != nil {
			return nil
		}
		msg.NameServers = append(msg.NameServers, ns)
	}
	// Parses additional records.
	for n := 0; n < int(msg.AR()); n++ {
		a, err := NewResourceWithReader(reader)
		if err != nil {
			return nil
		}
		msg.Additions = append(msg.Additions, a)
	}
	return nil
}

// Equals returns true if the message is same as the specified message, otherwise false.
func (msg *Message) Equals(other *Message) bool {
	if other == nil {
		return false
	}
	return bytes.Equal(msg.Bytes(), other.Bytes())
}

// Copy returns the copy message instance.
func (msg *Message) Copy() *Message {
	return &Message{
		Header:      NewHeaderWithBytes(msg.Header.bytes),
		Questions:   msg.Questions,
		Answers:     msg.Answers,
		NameServers: msg.NameServers,
		Additions:   msg.Additions,
	}
}

// Bytes returns the binary representation.
func (msg *Message) Bytes() []byte {
	bytes := msg.Header.Bytes()
	for _, q := range msg.Questions {
		bytes = append(bytes, q.Bytes()...)
	}
	for _, an := range msg.Answers {
		bytes = append(bytes, an.Bytes()...)
	}
	for _, ns := range msg.NameServers {
		bytes = append(bytes, ns.Bytes()...)
	}
	for _, a := range msg.Additions {
		bytes = append(bytes, a.Bytes()...)
	}
	return bytes
}

// String returns the string representation.
func (msg *Message) String() string {
	if msg == nil {
		return ""
	}
	return hex.EncodeToString(msg.Header.bytes)
}
