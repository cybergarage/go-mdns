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
	"bytes"
	"testing"
)

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
	})

	t.Run("Parse", func(t *testing.T) {
		testMsgs := [][]byte{
			{0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01},
		}

		header := NewRequestHeader()
		for _, testMsg := range testMsgs {
			if err := header.Parse(bytes.NewReader(testMsg)); err != nil {
				t.Error(err)
			}
		}
	})
}
