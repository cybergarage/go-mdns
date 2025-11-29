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

package mdnstest

import (
	"context"
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns"
)

func TestClient(t *testing.T) {
	log.EnableStdoutDebug(true)
	defer log.EnableStdoutDebug(false)

	client := mdns.NewClient()

	err := client.Start()
	if err != nil {
		t.Error(err)
		return
	}

	query := mdns.NewQuery(
		mdns.WithQueryServices(mdns.AutomaticBrowsingService),
	)

	_, err = client.Query(context.Background(), query)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}
