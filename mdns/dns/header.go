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

// RFC 1035: Domain names - implementation and specification
// https://www.rfc-editor.org/rfc/rfc1035.html
// RFC 6762: Multicast DNS
// https://www.rfc-editor.org/rfc/rfc6762.html

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/cybergarage/go-mdns/mdns/encoding"
)

const (
	headerSize = 12
)

type QR uint

const (
	Query    QR = 0
	Response QR = 1
)

type Opcode int

const (
	OpQuery  Opcode = 0
	OpIQuery Opcode = 1
	OpStatus Opcode = 2
)

type ResponseCode byte

const (
	NoError        ResponseCode = 0
	FormatError    ResponseCode = 1
	ServerFailure  ResponseCode = 2
	NameError      ResponseCode = 3
	NotImplemented ResponseCode = 4
	Refused        ResponseCode = 5
)

// Header represents a protocol header.
type Header struct {
	bytes []byte
}

// NewHeader returns a nil header instance.
func NewHeader() *Header {
	header := &Header{
		bytes: nil,
	}
	return header
}

// NewRequestHeader returns a request header instance.
func NewRequestHeader() *Header {
	header := &Header{
		bytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	return header
}

// NewResponseHeader returns a response header instance.
func NewResponseHeader() *Header {
	header := &Header{
		bytes: []byte{0x00, 0x00, 0x84, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
	return header
}

// NewHeaderWithBytes returns a header instance with the specified bytes.
func NewHeaderWithBytes(bytes []byte) *Header {
	return &Header{
		bytes: bytes,
	}
}

// NewHeaderWithReader returns a header instance with the specified reader.
func NewHeaderWithReader(reader io.Reader) (*Header, error) {
	header := &Header{
		bytes: nil,
	}
	return header, header.Parse(reader)
}

// Parse parses the specified reader.
func (header *Header) Parse(reader io.Reader) error {
	header.bytes = make([]byte, headerSize)
	_, err := reader.Read(header.bytes)
	return err
}

// Flags returns the flags.
func (header *Header) Flags() []byte {
	return header.bytes[2:4]
}

// ID returns the query identifier.
// RFC 6762: 18.1. ID (Query Identifier)
// In multicast query messages, the Query Identifier SHOULD be set to zero on transmission.
// In multicast responses, including unsolicited multicast responses, the Query Identifier MUST be set to zero on transmission, and MUST be ignored on reception.
func (header *Header) ID() uint {
	return encoding.BytesToInteger(header.bytes[:2])
}

// QR returns the query type.
// RFC 6762: 18.2. QR (Query/Response) Bit
// In query messages the QR bit MUST be zero. In response messages the QR bit MUST be one.
func (header *Header) QR() QR {
	if (header.bytes[2] & 0x80) == 0 {
		return Query
	}
	return Response
}

// Opcode returns the kind of query.
// RFC 6762: 18.3. OPCODE
// In both multicast query and multicast response messages, the OPCODE MUST be zero on transmission (only standard queries are currently supported over multicast).
func (header *Header) Opcode() Opcode {
	return Opcode((header.bytes[2] & 0x78) >> 3)
}

// AA returns the authoritative answer bit.
// RFC 6762: 18.4. AA (Authoritative Answer) Bit
// In query messages, the Authoritative Answer bit MUST be zero on transmission, and MUST be ignored on reception.
// In response messages for Multicast domains, the Authoritative Answer bit MUST be set to one (not setting this bit would imply there's some other place where "better" information may be found) and MUST be ignored on reception.
func (header *Header) AA() bool {
	return (header.bytes[2] & 0x04) == 0x04
}

// TC returns the truncated bit.
// RFC 6762: 18.5. TC (Truncated) Bit
// In query messages, if the TC bit is set, it means that additional Known-Answer records may be following shortly. A responder SHOULD record this fact, and wait for those additional Known-Answer records, before deciding whether to respond. If the TC bit is clear, it means that the querying host has no additional Known Answers.
// In multicast response messages, the TC bit MUST be zero on transmission, and MUST be ignored on reception.
func (header *Header) TC() bool {
	return (header.bytes[2] & 0x02) == 0x02
}

// RD returns the recursion desired bit.
// RFC 6762: 18.6. RD (Recursion Desired) Bit
// In both multicast query and multicast response messages, the Recursion Desired bit SHOULD be zero on transmission, and MUST be ignored on reception.
func (header *Header) RD() bool {
	return (header.bytes[2] & 0x01) == 0x01
}

// RA returns the recursion available bit.
// RFC 6762: 18.7. RA (Recursion Available) Bit
// In both multicast query and multicast response messages, the Recursion Available bit MUST be zero on transmission, and MUST be ignored on reception.
func (header *Header) RA() bool {
	return (header.bytes[3] & 0x80) == 0x80
}

// Z returns the zero bit.
// RFC 6762: 18.8. Z (Zero) Bit
// In both query and response messages, the Zero bit MUST be zero on transmission, and MUST be ignored on reception.
func (header *Header) Z() bool {
	return (header.bytes[3] & 0x40) == 0x40
}

// AD returns the authentic data bit.
// RFC 6762: 18.9. AD (Authentic Data) Bit
// In both multicast query and multicast response messages, the Authentic Data bit [RFC2535] MUST be zero on transmission, and MUST be ignored on reception.
func (header *Header) AD() bool {
	return (header.bytes[3] & 0x20) == 0x20
}

// CD returns the checking disabled bit.
// RFC 6762: 18.10. CD (Checking Disabled) Bit
// In both multicast query and multicast response messages, the Checking Disabled bit [RFC2535] MUST be zero on transmission, and MUST be ignored on reception.
func (header *Header) CD() bool {
	return (header.bytes[3] & 0x10) == 0x10
}

// ResponseCode returns the checking disabled bit.
// RFC 6762: 18.11. RCODE (Response Code)
// In both multicast query and multicast response messages, the Response Code MUST be zero on transmission. Multicast DNS messages received with non-zero Response Codes MUST be silently ignored.
func (header *Header) ResponseCode() ResponseCode {
	return ResponseCode(header.bytes[3] & 0x0F)
}

// setNumberOfEntries sets the specified number to the specified offset field.
func (header *Header) setNumberOfEntries(n uint, offset int) {
	bytes := make([]byte, 2)
	encoding.IntegerToBytes(n, bytes)
	header.bytes[offset] = bytes[0]
	header.bytes[offset+1] = bytes[1]
}

// setQD sets the specified number to the QD field.
func (header *Header) setQD(n uint) {
	header.setNumberOfEntries(n, 4)
}

// QD returns the number of entries in the question section.
func (header *Header) QD() uint {
	return encoding.BytesToInteger(header.bytes[4:6])
}

// setAN sets the specified number to the AN field.
func (header *Header) setAN(n uint) {
	header.setNumberOfEntries(n, 6)
}

// AN returns the number of entries in the answer section.
func (header *Header) AN() uint {
	return encoding.BytesToInteger(header.bytes[6:8])
}

// setNS sets the specified number to the NS field.
func (header *Header) setNS(n uint) {
	header.setNumberOfEntries(n, 8)
}

// NS returns the number of name server resource records in the authority records section.
func (header *Header) NS() uint {
	return encoding.BytesToInteger(header.bytes[8:10])
}

// setAR sets the specified number to the AR field.
func (header *Header) setAR(n uint) {
	header.setNumberOfEntries(n, 10)
}

// AR returns the number of resource records in the additional records section.
func (header *Header) AR() uint {
	return encoding.BytesToInteger(header.bytes[10:12])
}

// IsQuery returns true the QR bit is zero, otherwise false.
func (header *Header) IsQuery() bool {
	return header.QR() == Query
}

// IsResponse returns true the QR bit is one, otherwise false.
func (header *Header) IsResponse() bool {
	return (header.QR() == Response) && header.AA()
}

// Equal returns true if the header is same as the specified header, otherwise false.
func (header *Header) Equal(other *Header) bool {
	if other == nil {
		return false
	}
	return bytes.Equal(header.bytes, other.bytes)
}

// Copy returns the copy header instance.
func (header *Header) Copy() *Header {
	return NewHeaderWithBytes(header.bytes)
}

// Bytes returns the binary representation.
func (header *Header) Bytes() []byte {
	return header.bytes
}

// String returns the string representation.
func (header *Header) String() string {
	return hex.EncodeToString(header.bytes)
}
