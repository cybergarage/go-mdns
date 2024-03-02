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

// Question represents a question.
type Question struct {
	*record
}

// Questions represents a question array.
type Questions []*Question

// NewQuestion returns a new question innstance.
func NewQuestion() *Question {
	return &Question{
		record: newResourceRecord(),
	}
}

// NewQuestionWithRecord returns a new question innstance with the specified record.
func NewQuestionWithRecord(record *record) *Question {
	return &Question{
		record: record,
	}
}

// NewQuestionWithRecord returns a new question innstance with the specified record.
func NewQuestionWithReader(reader *Reader) (*Question, error) {
	r, err := newRequestRecordWithReader(reader)
	return NewQuestionWithRecord(r), err
}

// SetName sets the specified name to the question instance.
func (q *Question) SetName(name string) *Question {
	q.record.SetName(name)
	return q
}

// SetType sets the specified type to the question instance.
func (q *Question) SetType(t Type) *Question {
	q.record.SetType(t)
	return q
}

// SetClass sets the specified class to the question instance.
func (q *Question) SetClass(cls Class) *Question {
	q.record.SetClass(cls)
	return q
}

// SetUnicastResponse sets the specified unicast response flag.
func (q *Question) SetUnicastResponse(enabled bool) *Question {
	q.record.SetUnicastResponse(enabled)
	return q
}
