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
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/cybergarage/go-mdns/mdns/dns"
)

// The query message is built from the query options.
func TestQueryMessage(t *testing.T) {
	tests := []struct {
		name            string
		query           Query
		expectedName    string
		expectedType    Type
		expectedUnicast bool
	}{
		{
			name: "default",
			query: NewQuery(
				WithQueryService("_matterc._udp"),
			),
			// RFC 6763 browses a service type with a PTR query, and
			// RFC 6762 (5.4) does not set the unicast response bit by
			// default so that the other nodes can update their caches.
			expectedName:    "_matterc._udp.local",
			expectedType:    PTR,
			expectedUnicast: false,
		},
		{
			name: "subtype",
			query: NewQuery(
				WithQuerySubtype("_S3"),
				WithQueryService("_matterc._udp"),
			),
			expectedName:    "_S3._sub._matterc._udp.local",
			expectedType:    PTR,
			expectedUnicast: false,
		},
		{
			name: "instance",
			query: NewQuery(
				WithQueryName("DD200C20D25AE5F7._matterc._udp.local"),
				WithQueryType(SRV),
				WithQueryUnicastResponse(true),
			),
			expectedName:    "DD200C20D25AE5F7._matterc._udp.local",
			expectedType:    SRV,
			expectedUnicast: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.query.Name() != test.expectedName {
				t.Errorf("name %s != %s", test.query.Name(), test.expectedName)
			}

			msg := NewRequestWithQuery(test.query)
			questions := msg.Questions()
			if len(questions) != 1 {
				t.Fatalf("question count %d != 1", len(questions))
			}

			q := questions[0]
			if !q.IsName(test.expectedName) {
				t.Errorf("question name %s != %s", q.Name(), test.expectedName)
			}
			if q.Type() != test.expectedType {
				t.Errorf("question type %s != %s", q.Type(), test.expectedType)
			}
			if q.IsUnicastResponse() != test.expectedUnicast {
				t.Errorf("unicast response %t != %t", q.IsUnicastResponse(), test.expectedUnicast)
			}
			if !q.Class().Equal(dns.IN) {
				t.Errorf("question class %d is not IN", q.Class())
			}
		})
	}
}

// nolint: gocyclo
func TestQuery(t *testing.T) {
	tests := []struct {
		name  string
		query Query
	}{
		{
			name: AutomaticBrowsingService,
			query: NewQuery(
				WithQueryService(DefaultQueryService),
				WithQueryDomain(DefaultQueryDomain),
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			qmsg := NewRequestWithQuery(test.query)
			qmsgBytes := qmsg.Bytes()
			msg, err := dns.NewMessageWithBytes(qmsgBytes)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(qmsg.Bytes(), msg.Bytes()) {
				t.Errorf("%s != %s", hex.EncodeToString(qmsg.Bytes()), hex.EncodeToString(msg.Bytes()))
			}
		})
	}
}
