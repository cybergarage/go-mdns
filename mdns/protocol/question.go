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
	"io"

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

// Question represents a question.
type Question struct {
	Name            string
	Type            Type
	UnicastResponse bool
	Class           Class
}

// Questions represents a question array.
type Questions []*Question

// NewQuestion returns a new question innstance.
func NewQuestion() *Question {
	return &Question{
		Name:            "",
		Type:            0,
		UnicastResponse: false,
		Class:           IN,
	}
}

// NewQuestionWithReader returns a new question innstance with the specified reader.
func NewQuestionWithReader(reader io.Reader) (*Question, error) {
	q := NewQuestion()
	return q, q.Parse(reader)
}

// SetName sets the specified name to the question instance.
func (q *Question) SetName(name string) *Question {
	q.Name = name
	return q
}

// Parse parses the specified reader.
func (q *Question) Parse(reader io.Reader) error {
	var err error

	// Parses domain names
	q.Name, err = parseName(reader)
	if err != nil {
		return err
	}

	// Parses query type
	typeBytes := make([]byte, 2)
	_, err = reader.Read(typeBytes)
	if err != nil {
		return err
	}
	q.Type = Type(encoding.BytesToInteger(typeBytes))

	// Parses class type
	classBytes := make([]byte, 2)
	_, err = reader.Read(classBytes)
	if err != nil {
		return err
	}
	class := encoding.BytesToInteger(classBytes)
	q.UnicastResponse = false
	if (class & unicastResponseMask) != 0 {
		q.UnicastResponse = true
	}
	q.Class = Class(class & classMask)

	return nil
}

// Bytes returns the binary representation.
func (q *Question) Bytes() []byte {
	bytes := nameToBytes(q.Name)

	typeBytes := make([]byte, 2)
	bytes = append(bytes, encoding.IntegerToBytes(uint(q.Type), typeBytes)...)

	classBytes := make([]byte, 2)
	class := q.Class
	if q.UnicastResponse {
		class |= unicastResponseMask
	}
	bytes = append(bytes, encoding.IntegerToBytes(uint(class), classBytes)...)

	return bytes
}
