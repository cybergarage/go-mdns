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

// question represents a question.
type question struct {
	*record
}

// QuestionOption represents a question option.
type QuestionOption func(*question)

// Questions represents a question array.
type Questions []Question

// WithQuestionName sets the question name.
func WithQuestionName(name string) QuestionOption {
	return func(q *question) {
		q.SetName(name)
	}
}

// WithQuestionType sets the question type.
func WithQuestionType(t Type) QuestionOption {
	return func(q *question) {
		q.SetType(t)
	}
}

// WithQuestionClass sets the question class.
func WithQuestionClass(cls Class) QuestionOption {
	return func(q *question) {
		q.SetClass(cls)
	}
}

// NewQuestion returns a new question instance.
func NewQuestion(opts ...QuestionOption) Question {
	q := &question{
		record: newRecord(),
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// NewQuestionWithRecord returns a new question instance with the specified record.
func NewQuestionWithRecord(record *record) Question {
	return &question{
		record: record,
	}
}

// NewQuestionWithReader returns a new question instance with the specified record.
func NewQuestionWithReader(reader *Reader) (Question, error) {
	r, err := NewRequestRecordWithReader(reader)
	return NewQuestionWithRecord(r), err
}

// IsUnicastResponse returns true if the question has the unicast response bit set, otherwise false.
func (q *question) IsUnicastResponse() bool {
	return (q.Class() & QU) == QU
}

// Equal returns true if this record is equal to  the specified resource record. otherwise false.
func (q *question) Equal(other Record) bool {
	return EqualContent(q, other)
}
