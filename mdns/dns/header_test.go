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
	"bytes"
	"testing"
)

// nolint: gocyclo
func TestHeader(t *testing.T) {
	t.Run("RequestHeader", func(t *testing.T) {
		header := NewRequestHeader()
		if header.ID() != 0 {
			t.Errorf("%d != %d", header.ID(), 0)
		}
		if header.QR() != Query {
			t.Errorf("%d != %d", header.QR(), Query)
		}
		if header.Opcode() != OpQuery {
			t.Errorf("%d != %d", header.Opcode(), OpQuery)
		}
		if header.AA() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.TC() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.RD() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.RA() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.Z() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.AD() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.CD() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.ResponseCode() != NoError {
			t.Errorf("%b != %b", header.ResponseCode(), NoError)
		}
		if header.QD() != 0 {
			t.Errorf("%d != %d", header.QD(), 0)
		}
		if header.AN() != 0 {
			t.Errorf("%d != %d", header.AN(), 0)
		}
		if header.NS() != 0 {
			t.Errorf("%d != %d", header.NS(), 0)
		}
		if header.AR() != 0 {
			t.Errorf("%d != %d", header.AR(), 0)
		}
		if !header.IsQuery() {
			t.Errorf("%b", header.QR())
		}
	})

	t.Run("ResponseHeader", func(t *testing.T) {
		header := NewResponseHeader()
		if header.ID() != 0 {
			t.Errorf("%d != %d", header.ID(), 0)
		}
		if header.QR() != Response {
			t.Errorf("%d != %d", header.QR(), Response)
		}
		if header.Opcode() != OpQuery {
			t.Errorf("%d != %d", header.Opcode(), OpQuery)
		}
		if !header.AA() {
			t.Errorf("%t != %t", header.AA(), true)
		}
		if header.TC() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.RD() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.RA() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.Z() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.AD() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.CD() {
			t.Errorf("%t != %t", header.AA(), false)
		}
		if header.ResponseCode() != NoError {
			t.Errorf("%b != %b", header.ResponseCode(), NoError)
		}
		if header.QD() != 0 {
			t.Errorf("%d != %d", header.QD(), 0)
		}
		if header.AN() != 0 {
			t.Errorf("%d != %d", header.AN(), 0)
		}
		if header.NS() != 0 {
			t.Errorf("%d != %d", header.NS(), 0)
		}
		if header.AR() != 0 {
			t.Errorf("%d != %d", header.AR(), 0)
		}
		if !header.IsResponse() {
			t.Errorf("%b", header.QR())
		}
	})

	t.Run("ParseRequest", func(t *testing.T) {
		testMsgs := [][]byte{
			{0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01},
		}

		for _, testMsg := range testMsgs {
			header, err := NewHeaderWithReader(bytes.NewReader(testMsg))
			if err != nil {
				t.Error(err)
			}
			if !header.IsQuery() {
				t.Errorf("%b", header.QR())
			}
		}
	})

	t.Run("ParseResponse", func(t *testing.T) {
		testMsgs := [][]byte{
			{0x00, 0x00, 0x84, 0x00, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x04},
		}

		for _, testMsg := range testMsgs {
			header, err := NewHeaderWithReader(bytes.NewReader(testMsg))
			if err != nil {
				t.Error(err)
			}
			if !header.IsResponse() {
				t.Errorf("%b", header.QR())
			}
		}
	})
}
