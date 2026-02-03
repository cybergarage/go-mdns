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
	"regexp"
)

// MessageOption represents a message option.
type MessageOption func(*message) error

// Message represents a DNS message.
type Message interface {
	// From returns the source address of the message.
	From() Addr
	// Flags returns the flags.
	Flags() []byte
	// ID returns the query identifier.
	// RFC 6762: 18.1. ID (Query Identifier)
	// In multicast query messages, the Query Identifier SHOULD be set to zero on transmission.
	// In multicast responses, including unsolicited multicast responses, the Query Identifier MUST be set to zero on transmission, and MUST be ignored on reception.
	ID() uint
	// QR returns the query type.
	// RFC 6762: 18.2. QR (Query/Response) Bit
	// In query messages the QR bit MUST be zero. In response messages the QR bit MUST be one.
	QR() QR
	// Opcode returns the kind of query.
	// RFC 6762: 18.3. OPCODE
	// In both multicast query and multicast response messages, the OPCODE MUST be zero on transmission (only standard queries are currently supported over multicast).
	Opcode() Opcode
	// AA returns the authoritative answer bit.
	// RFC 6762: 18.4. AA (Authoritative Answer) Bit
	// In query messages, the Authoritative Answer bit MUST be zero on transmission, and MUST be ignored on reception.
	// In response messages for Multicast domains, the Authoritative Answer bit MUST be set to one (not setting this bit would imply there's some other place where "better" information may be found) and MUST be ignored on reception.
	AA() bool
	// TC returns the truncated bit.
	// RFC 6762: 18.5. TC (Truncated) Bit
	// In query messages, if the TC bit is set, it means that additional Known-Answer records may be following shortly. A responder SHOULD record this fact, and wait for those additional Known-Answer records, before deciding whether to respond. If the TC bit is clear, it means that the querying host has no additional Known Answers.
	// In multicast response messages, the TC bit MUST be zero on transmission, and MUST be ignored on reception.
	TC() bool
	// RD returns the recursion desired bit.
	// RFC 6762: 18.6. RD (Recursion Desired) Bit
	// In both multicast query and multicast response messages, the Recursion Desired bit SHOULD be zero on transmission, and MUST be ignored on reception.
	RD() bool
	// RA returns the recursion available bit.
	// RFC 6762: 18.7. RA (Recursion Available) Bit
	// In both multicast query and multicast response messages, the Recursion Available bit MUST be zero on transmission, and MUST be ignored on reception.
	RA() bool
	// Z returns the zero bit.
	// RFC 6762: 18.8. Z (Zero) Bit
	// In both query and response messages, the Zero bit MUST be zero on transmission, and MUST be ignored on reception.
	Z() bool
	// AD returns the authentic data bit.
	// RFC 6762: 18.9. AD (Authentic Data) Bit
	// In both multicast query and multicast response messages, the Authentic Data bit [RFC2535] MUST be zero on transmission, and MUST be ignored on reception.
	AD() bool
	// CD returns the checking disabled bit.
	// RFC 6762: 18.10. CD (Checking Disabled) Bit
	// In both multicast query and multicast response messages, the Checking Disabled bit [RFC2535] MUST be zero on transmission, and MUST be ignored on reception.
	CD() bool
	// ResponseCode returns the checking disabled bit.
	// RFC 6762: 18.11. RCODE (Response Code)
	// In both multicast query and multicast response messages, the Response Code MUST be zero on transmission. Multicast DNS messages received with non-zero Response Codes MUST be silently ignored.
	ResponseCode() ResponseCode
	// QD returns the number of entries in the question section.
	QD() uint
	// AN returns the number of entries in the answer section.
	AN() uint
	// NS returns the number of name server resource records in the authority records section.
	NS() uint
	// AR returns the number of resource records in the additional records section.
	AR() uint
	// IsQuery returns true the QR bit is zero, otherwise false.
	IsQuery() bool
	// IsResponse returns true the QR bit is one, otherwise false.
	IsResponse() bool
	// Questions returns all questions in the message.
	Questions() Questions
	// Answers returns all answers in the message.
	Answers() ResourceRecordSet
	// NameServers returns all name servers in the message.
	NameServers() ResourceRecordSet
	// Additions returns all additional records in the message.
	Additions() ResourceRecordSet
	// RecordSet returns all records which includes questions, answers, name servers, and additions.
	RecordSet() RecordSet
	// ResourceRecordSet returns only all resource records in the message without questions.
	ResourceRecordSet() ResourceRecordSet
	// LookupResourceRecordByName returns the resource record of the specified name.
	LookupResourceRecordByName(name string) (ResourceRecord, bool)
	// LookupResourceRecordByNameRegex returns the resource record of the specified name regex.
	LookupResourceRecordByNameRegex(re *regexp.Regexp) (ResourceRecord, bool)
	// Equal returns true if the message is same as the specified message, otherwise false.
	Equal(other Message) bool
	// Copy returns the copy message instance.
	Copy() Message
	// Bytes returns the byte representation of the message.
	Bytes() []byte
	// String returns the string representation of the message.
	String() string
	// MessageHelper represents a message helper functions.
	MessageHelper
}

// MessageHelper represents a message helper functions.
type MessageHelper interface {
	// IsQueryWithUnicastResponse returns true if the message is a query with the unicast response bit set, otherwise false.
	IsQueryWithUnicastResponse() bool
	// IsQueryAnswer returns true if the message is a response to a query, otherwise false.
	IsQueryAnswer(msg Message) bool
	// LookupResourceRecordByNamePrefix returns the resource record of the specified name prefix.
	LookupResourceRecordByNamePrefix(prefix string) (ResourceRecord, bool)
	// LookupResourceRecordByNameSuffix returns the resource record of the specified name suffix.
	LookupResourceRecordByNameSuffix(suffix string) (ResourceRecord, bool)
}
