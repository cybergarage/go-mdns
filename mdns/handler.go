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
	"reflect"
	"sync"
)

type msgHandler struct {
	sync.Mutex
	handlers []MessageHandler
}

func newMessageHandler() *msgHandler {
	msg := &msgHandler{
		Mutex:    sync.Mutex{},
		handlers: []MessageHandler{},
	}
	return msg
}

// RegisterMessageHandler adds a message handler to the server.
func (msg *msgHandler) RegisterMessageHandler(handler MessageHandler) {
	msg.Lock()
	defer msg.Unlock()
	// Check if handler already exists using reflection
	handlerVal := reflect.ValueOf(handler)
	for _, h := range msg.handlers {
		if reflect.ValueOf(h).Pointer() == handlerVal.Pointer() {
			return
		}
	}
	msg.handlers = append(msg.handlers, handler)
}

// UnRegisterMessageHandler removes a message handler from the server.
func (msg *msgHandler) UnRegisterMessageHandler(handler MessageHandler) {
	msg.Lock()
	defer msg.Unlock()
	// Use reflection to compare function pointers
	handlerVal := reflect.ValueOf(handler)
	for i, h := range msg.handlers {
		if reflect.ValueOf(h).Pointer() == handlerVal.Pointer() {
			msg.handlers = append(msg.handlers[:i], msg.handlers[i+1:]...)
			return
		}
	}
}

// processMessageHandlers processes the message using registered handlers.
func (msg *msgHandler) processMessageHandlers(message Message) {
	for _, handler := range msg.handlers {
		handler(message)
	}
}
