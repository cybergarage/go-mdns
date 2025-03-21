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

package dns

import (
	"bytes"
	"fmt"
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

// NewMessageWithBytes returns a message instance with the specified bytes.
func NewMessageWithBytes(msgBytes []byte) (*Message, error) {
	msg := NewMessage()
	if err := msg.Parse(msgBytes); err != nil {
		return nil, err
	}
	return msg, nil
}

// AddQuestion adds the specified question into the message.
func (msg *Message) AddQuestion(q *Question) {
	msg.Questions = append(msg.Questions, q)
	msg.setQD(uint(len(msg.Questions)))
}

// AddAnswer adds the specified answer into the message.
func (msg *Message) AddAnswer(a Answer) {
	msg.Answers = append(msg.Answers, a)
	msg.setAN(uint(len(msg.Answers)))
}

// AddNameServer adds the specified name server into the message.
func (msg *Message) AddNameServer(ns NameServer) {
	msg.NameServers = append(msg.NameServers, ns)
	msg.setNS(uint(len(msg.NameServers)))
}

// AddAddition adds the specified additional record into the message.
func (msg *Message) AddAddition(a Addition) {
	msg.Additions = append(msg.Additions, a)
	msg.setAR(uint(len(msg.Additions)))
}

// Parse parses the specified reader.
func (msg *Message) Parse(msgBytes []byte) error {
	reader := NewReaderWithBytes(msgBytes)
	if err := msg.Header.Parse(reader); err != nil {
		return fmt.Errorf("header : %w", err)
	}

	// Parses questions.
	for n := 0; n < int(msg.QD()); n++ {
		r, err := NewRequestRecordWithReader(reader)
		if err != nil {
			return fmt.Errorf("question[%d] : %w", n, err)
		}
		msg.Questions = append(msg.Questions, NewQuestionWithRecord(r))
	}
	// Parses answers.
	for n := 0; n < int(msg.AN()); n++ {
		a, err := NewResourceRecordWithReader(reader)
		if err != nil {
			return fmt.Errorf("answer[%d] : %w", n, err)
		}
		msg.Answers = append(msg.Answers, a)
	}
	// Parses authorities.
	for n := 0; n < int(msg.NS()); n++ {
		ns, err := NewResourceRecordWithReader(reader)
		if err != nil {
			return fmt.Errorf("authority[%d] : %w", n, err)
		}
		msg.NameServers = append(msg.NameServers, ns)
	}
	// Parses additional records.
	for n := 0; n < int(msg.AR()); n++ {
		a, err := NewResourceRecordWithReader(reader)
		if err != nil {
			return fmt.Errorf("additional[%d] : %w", n, err)
		}
		msg.Additions = append(msg.Additions, a)
	}
	return nil
}

// Equal returns true if the message is same as the specified message, otherwise false.
func (msg *Message) Equal(other *Message) bool {
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

// Records returns all resource records.
func (msg *Message) Records() Records {
	records := Records{}
	for _, r := range msg.Questions {
		records = append(records, r)
	}
	records = append(records, msg.Answers...)
	records = append(records, msg.Answers...)
	records = append(records, msg.NameServers...)
	records = append(records, msg.Additions...)
	return records
}

// ResourceRecords returns all resource records.
func (msg *Message) ResourceRecords() ResourceRecords {
	records := ResourceRecords{}
	records = append(records, msg.Answers...)
	records = append(records, msg.NameServers...)
	records = append(records, msg.Additions...)
	return records
}

// LookupResourceRecordForName returns the resource record of the specified name.
func (msg *Message) LookupResourceRecordForName(name string) (ResourceRecord, bool) {
	return msg.ResourceRecords().LookupRecordForName(name)
}

// LookupResourceRecordForNamePrefix returns the resource record of the specified name prefix.
func (msg *Message) LookupResourceRecordForNamePrefix(prefix string) (ResourceRecord, bool) {
	return msg.ResourceRecords().LookupRecordForNamePrefix(prefix)
}

// LookupResourceRecordForNameSuffix returns the resource record of the specified name suffix.
func (msg *Message) LookupResourceRecordForNameSuffix(prefix string) (ResourceRecord, bool) {
	return msg.ResourceRecords().LookupRecordForNameSuffix(prefix)
}

// HasResourceRecord returns true if the resource record of the specified name is included in the message. otherwise false.
func (msg *Message) HasResourceRecord(name string) bool {
	_, ok := msg.LookupResourceRecordForName(name)
	return ok
}

// String returns the string representation.
func (msg *Message) String() string {
	return msg.Records().String()
}

// Bytes returns the binary representation.
func (msg *Message) Bytes() []byte {
	bytes := msg.Header.Bytes()
	for _, q := range msg.Questions {
		if b, err := q.RequestBytes(); err == nil {
			bytes = append(bytes, b...)
		}
	}
	for _, an := range msg.Answers {
		if b, err := an.ResponseBytes(); err == nil {
			bytes = append(bytes, b...)
		}
	}
	for _, ns := range msg.NameServers {
		if b, err := ns.ResponseBytes(); err == nil {
			bytes = append(bytes, b...)
		}
	}
	for _, a := range msg.Additions {
		if b, err := a.ResponseBytes(); err == nil {
			bytes = append(bytes, b...)
		}
	}
	return bytes
}
