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
	"net"
	"regexp"
)

// message represents a protocol message.
type message struct {
	*Header
	from        net.Addr
	questions   Questions
	answers     Answers
	nameServers NameServers
	additions   Additions
}

func newMessage(opts ...MessageOption) *message {
	msg := &message{
		Header:      NewHeader(),
		from:        nil,
		questions:   Questions{},
		answers:     Answers{},
		nameServers: NameServers{},
		additions:   Additions{},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// WithMessageQuestions returns a message option with the specified questions.
func WithMessageQuestions(questions ...Question) MessageOption {
	return func(msg *message) error {
		for _, q := range questions {
			msg.AddQuestion(q)
		}
		return nil
	}
}

// WithMessageFrom returns a message option with the specified source address.
func WithMessageFrom(addr net.Addr) MessageOption {
	return func(msg *message) error {
		msg.from = addr
		return nil
	}
}

// NewMessage returns a nil message instance.
func NewMessage(opts ...MessageOption) Message {
	return newMessage(opts...)
}

// NewRequestMessage returns a request message instance.
func NewRequestMessage(opts ...MessageOption) Message {
	msg := newMessage()
	msg.Header = NewRequestHeader()
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// NewResponseMessage returns a response message instance.
func NewResponseMessage(opts ...MessageOption) Message {
	msg := newMessage()
	msg.Header = NewResponseHeader()
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// NewMessageWithBytes returns a message instance with the specified bytes.
func NewMessageWithBytes(msgBytes []byte, opts ...MessageOption) (Message, error) {
	msg := newMessage()
	if err := msg.Parse(msgBytes); err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg, nil
}

// From returns the source address of the message.
func (msg *message) From() net.Addr {
	return msg.from
}

// AddQuestion adds the specified question into the message.
func (msg *message) AddQuestion(q Question) {
	msg.questions = append(msg.questions, q)
	msg.setQD(uint(len(msg.questions)))
}

// AddAnswer adds the specified answer into the message.
func (msg *message) AddAnswer(a Answer) {
	msg.answers = append(msg.answers, a)
	msg.setAN(uint(len(msg.answers)))
}

// AddNameServer adds the specified name server into the message.
func (msg *message) AddNameServer(ns NameServer) {
	msg.nameServers = append(msg.nameServers, ns)
	msg.setNS(uint(len(msg.nameServers)))
}

// AddAddition adds the specified additional record into the message.
func (msg *message) AddAddition(a Addition) {
	msg.additions = append(msg.additions, a)
	msg.setAR(uint(len(msg.additions)))
}

// Parse parses the specified reader.
func (msg *message) Parse(msgBytes []byte) error {
	reader := NewReaderWithBytes(msgBytes)
	if err := msg.Header.Parse(reader); err != nil {
		return fmt.Errorf("header : %w", err)
	}

	// Parses questions.
	for n := range int(msg.QD()) {
		r, err := NewRequestRecordWithReader(reader)
		if err != nil {
			return fmt.Errorf("question[%d] : %w", n, err)
		}
		msg.questions = append(msg.questions, NewQuestionWithRecord(r))
	}
	// Parses answers.
	for n := range int(msg.AN()) {
		a, err := NewResourceRecordWithReader(reader)
		if err != nil {
			return fmt.Errorf("answer[%d] : %w", n, err)
		}
		msg.answers = append(msg.answers, a)
	}
	// Parses authorities.
	for n := range int(msg.NS()) {
		ns, err := NewResourceRecordWithReader(reader)
		if err != nil {
			return fmt.Errorf("authority[%d] : %w", n, err)
		}
		msg.nameServers = append(msg.nameServers, ns)
	}
	// Parses additional records.
	for n := range int(msg.AR()) {
		a, err := NewResourceRecordWithReader(reader)
		if err != nil {
			return fmt.Errorf("additional[%d] : %w", n, err)
		}
		msg.additions = append(msg.additions, a)
	}
	return nil
}

// Equal returns true if the message is same as the specified message, otherwise false.
func (msg *message) Equal(other *message) bool {
	if other == nil {
		return false
	}
	return bytes.Equal(msg.Bytes(), other.Bytes())
}

// Copy returns the copy message instance.
func (msg *message) Copy() Message {
	return &message{
		Header:      NewHeaderWithBytes(msg.Header.bytes),
		from:        msg.from,
		questions:   msg.questions,
		answers:     msg.answers,
		nameServers: msg.nameServers,
		additions:   msg.additions,
	}
}

// Questions returns all questions in the message.
func (msg *message) Questions() Questions {
	return msg.questions
}

// Answers returns all answers in the message.
func (msg *message) Answers() ResourceRecordSet {
	return msg.answers
}

// NameServers returns all name servers in the message.
func (msg *message) NameServers() ResourceRecordSet {
	return msg.nameServers
}

// Additions returns all additional records in the message.
func (msg *message) Additions() ResourceRecordSet {
	return msg.additions
}

// RecordSet returns all records which includes questions, answers, name servers, and additions.
func (msg *message) RecordSet() RecordSet {
	records := RecordSet{}
	for _, r := range msg.questions {
		records = append(records, r)
	}
	records = append(records, msg.answers...)
	records = append(records, msg.nameServers...)
	records = append(records, msg.additions...)
	return records
}

// ResourceRecordSet returns only all resource records in the message without questions.
func (msg *message) ResourceRecordSet() ResourceRecordSet {
	records := ResourceRecordSet{}
	records = append(records, msg.answers...)
	records = append(records, msg.nameServers...)
	records = append(records, msg.additions...)
	return records
}

// LookupResourceRecordByName returns the resource record of the specified name.
func (msg *message) LookupResourceRecordByName(name string) (ResourceRecord, bool) {
	return msg.ResourceRecordSet().LookupRecordByName(name)
}

// LookupResourceRecordByNameRegex returns the resource record of the specified name regex.
func (msg *message) LookupResourceRecordByNameRegex(re *regexp.Regexp) (ResourceRecord, bool) {
	return msg.ResourceRecordSet().LookupRecordByNameRegex(re)
}

// LookupResourceRecordByNamePrefix returns the resource record of the specified name prefix.
func (msg *message) LookupResourceRecordByNamePrefix(prefix string) (ResourceRecord, bool) {
	return msg.ResourceRecordSet().LookupRecordByNamePrefix(prefix)
}

// LookupResourceRecordByNameSuffix returns the resource record of the specified name suffix.
func (msg *message) LookupResourceRecordByNameSuffix(suffix string) (ResourceRecord, bool) {
	return msg.ResourceRecordSet().LookupRecordByNameSuffix(suffix)
}

// String returns the string representation.
func (msg *message) String() string {
	return msg.RecordSet().String()
}

// Bytes returns the binary representation.
func (msg *message) Bytes() []byte {
	bytes := msg.Header.Bytes()
	for _, q := range msg.questions {
		if b, err := q.RequestBytes(); err == nil {
			bytes = append(bytes, b...)
		}
	}
	for _, an := range msg.answers {
		if b, err := an.ResponseBytes(); err == nil {
			bytes = append(bytes, b...)
		}
	}
	for _, ns := range msg.nameServers {
		if b, err := ns.ResponseBytes(); err == nil {
			bytes = append(bytes, b...)
		}
	}
	for _, a := range msg.additions {
		if b, err := a.ResponseBytes(); err == nil {
			bytes = append(bytes, b...)
		}
	}
	return bytes
}
