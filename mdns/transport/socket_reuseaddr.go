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

//go:build linux || windows

package transport

import (
	"os"
	"syscall"
)

// SetReuseAddrFd sets a flag to SO_REUSEADDR.
// nolint: nosnakecase
func (sock *Socket) SetReuseAddrFd(fd uintptr, flag bool) error {
	opt := 0
	if flag {
		opt = 1
	}
	return syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, opt)
}

// SetReuseAddr sets a flag to SO_REUSEADDR and SO_REUSEPORT.
// nolint: nosnakecase
func (sock *Socket) SetReuseAddr(file *os.File, flag bool) error {
	return sock.SetReuseAddrFd(file.Fd(), flag)
}
