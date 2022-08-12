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
		Class:           0,
	}
}

// NewQuestionWithReader returns a new question innstance with the specified reader.
func NewQuestionWithReader(reader io.Reader) (*Question, error) {
	q := NewQuestion()
	return q, q.Parse(reader)
}

// Parse parses the specified reader.
func (q *Question) Parse(reader io.Reader) error {
	// Parses domain names
	nextNameLenBuf := make([]byte, 1)
	_, err := reader.Read(nextNameLenBuf)
	for err == nil {
		nextNameLen := encoding.BytesToInteger(nextNameLenBuf)
		if nextNameLen == 0 {
			break
		}
		nextName := make([]byte, nextNameLen)
		_, err = reader.Read(nextName)
		if err != nil {
			return err
		}
		if 0 < len(q.Name) {
			q.Name += "."
		}
		q.Name += string(nextName)
		_, err = reader.Read(nextNameLenBuf)
	}
	if err != nil {
		return err
	}

	// Parses query type
	queryTypeBuf := make([]byte, 2)
	_, err = reader.Read(queryTypeBuf)
	if err != nil {
		return err
	}
	q.Type = Type(encoding.BytesToInteger(queryTypeBuf))

	// Parses c;ass type
	classBuf := make([]byte, 2)
	_, err = reader.Read(classBuf)
	if err != nil {
		return err
	}
	class := encoding.BytesToInteger(classBuf)
	q.UnicastResponse = false
	if (class & unicastResponseMask) != 0 {
		q.UnicastResponse = true
	}
	q.Class = Class(class & classMask)

	return nil
}
