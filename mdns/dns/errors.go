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
	"errors"
	"fmt"
)

var (
	// ErrNil indicates that a required value is nil.
	ErrNil = errors.New("nil")
	// ErrInvalid indicates that a value is invalid, such as a malformed
	// message or a name which cannot be encoded.
	ErrInvalid = errors.New("invalid")
	// ErrNilReader indicates that a record has no reader to parse its data
	// with. It wraps ErrNil.
	ErrNilReader = fmt.Errorf("reader is %w", ErrNil)
)
