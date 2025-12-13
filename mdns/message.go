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
	"github.com/cybergarage/go-mdns/mdns/dns"
)

// Message represents a protocol message.
type Message = dns.Message

// NewRequestWithQuery returns a nil message instance.
func NewRequestWithQuery(query Query) Message {
	questions := []dns.Question{}
	for _, service := range query.Services() {
		q := dns.NewQuestion(
			dns.WithQuestionName(dns.NewNameWithStrings(service, query.Domain())),
			dns.WithQuestionType(dns.PTR),
			dns.WithQuestionClass(dns.IN),
		)
		questions = append(questions, q)
	}
	return dns.NewRequestMessage(dns.WithMessageQuestions(questions...))
}
