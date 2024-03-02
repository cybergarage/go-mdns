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

package mdns

import (
	"github.com/cybergarage/go-mdns/mdns/protocol"
)

// Message represents a protocol message.
type Message = protocol.Message

// NewMessage returns a nil message instance.
func NewRequestWithQuery(query *Query) *Message {
	msg := protocol.NewRequestMessage()
	q := protocol.NewQuestion()
	q.SetName(query.String())
	q.SetType(protocol.PTR)
	q.SetClass(protocol.IN)
	msg.AddQuestion(q)
	return msg
}

// NewRequestWithQueries returns a new request message with the specified queries.
func NewRequestWithQueries(queries []*Query) *Message {
	msg := protocol.NewRequestMessage()
	for _, query := range queries {
		q := protocol.NewQuestion()
		q.SetName(query.String())
		q.SetType(protocol.PTR)
		q.SetClass(protocol.IN)
		msg.AddQuestion(q)
	}
	return msg
}
