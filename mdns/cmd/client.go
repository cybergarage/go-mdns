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

package cmd

import (
	"fmt"
	"net"

	"github.com/cybergarage/go-mdns/mdns"
	"github.com/spf13/viper"
)

// newClient returns a new mDNS client which is configured by the common flags.
func newClient() (mdns.Client, error) {
	opts := []mdns.ClientOption{
		mdns.WithClientQueryTimeout(queryTimeout()),
	}

	ifnames := viper.GetStringSlice(InterfaceParamStr)
	if 0 < len(ifnames) {
		ifis := make([]*net.Interface, 0, len(ifnames))
		for _, ifname := range ifnames {
			ifi, err := net.InterfaceByName(ifname)
			if err != nil {
				return nil, err
			}
			ifis = append(ifis, ifi)
		}
		opts = append(opts, mdns.WithClientInterfaces(ifis...))
	}

	switch family := viper.GetString(FamilyParamStr); family {
	case FamilyAllStr, "":
	case FamilyIPv4Str:
		opts = append(opts, mdns.WithClientIPv6Enabled(false))
	case FamilyIPv6Str:
		opts = append(opts, mdns.WithClientIPv4Enabled(false))
	default:
		return nil, fmt.Errorf("invalid address family: %s", family)
	}

	return mdns.NewClient(opts...), nil
}

// startClient returns a started mDNS client and its stop function.
func startClient() (mdns.Client, func(), error) {
	client, err := newClient()
	if err != nil {
		return nil, nil, err
	}

	if err := client.Start(); err != nil {
		return nil, nil, err
	}

	return client, func() {
		client.Stop()
	}, nil
}
