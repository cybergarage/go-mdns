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

// Question represents a question.
type Question struct {
	*record
}

// Questions represents a question array.
type Questions []*Question

// NewQuestion returns a new question instance.
func NewQuestion() *Question {
	return &Question{
		record: newResourceRecord(),
	}
}

// NewQuestionWithRecord returns a new question instance with the specified record.
func NewQuestionWithRecord(record *record) *Question {
	return &Question{
		record: record,
	}
}

// NewQuestionWithRecord returns a new question instance with the specified record.
func NewQuestionWithReader(reader *Reader) (*Question, error) {
	r, err := NewRequestRecordWithReader(reader)
	return NewQuestionWithRecord(r), err
}
