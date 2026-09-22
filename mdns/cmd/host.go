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
	"context"

	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   "host <hostname>",
	Short: "Resolve a host name to its addresses",
	Long: `Resolve a host name to its addresses with Multicast DNS (RFC 6762).

The IPv6 scoped addressing zone is set for the link-local addresses, so that the
printed addresses can be used to connect to the host.`,
	Example: `  mdnslookup host macmini.local
  mdnslookup host --family ipv6 84FCE6036F38.local`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, stop, err := startClient()
		if err != nil {
			return err
		}
		defer stop()

		ctx, cancel := context.WithTimeout(cmd.Context(), queryTimeout())
		defer cancel()

		host := args[0]
		addrs, err := client.LookupHost(ctx, host)
		if err != nil {
			return err
		}

		return outputHostAddrs(host, addrs)
	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
}
