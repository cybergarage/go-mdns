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
)

// Question represents a question.
type Question struct {
	*Record
}

// Questions represents a question array.
type Questions []*Question

// NewQuestion returns a new question innstance.
func NewQuestion() *Question {
	return &Question{
		Record: newResourceRecord(),
	}
}

// NewQuestionWithReader returns a new question innstance with the specified reader.
func NewQuestionWithReader(reader io.Reader) (*Question, error) {
	q := NewQuestion()
	return q, q.Parse(reader)
}

// SetName sets the specified name to the question instance.
func (q *Question) SetName(name string) *Question {
	q.Record.SetName(name)
	return q
}

// SetType sets the specified type to the question instance.
func (q *Question) SetType(t Type) *Question {
	q.Record.SetType(t)
	return q
}

// Parse parses the specified reader.
func (q *Question) Parse(reader io.Reader) error {
	return q.Record.Parse(reader)
}
