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
	"fmt"
	"strings"

	"github.com/cybergarage/go-mdns/mdns"
	"github.com/cybergarage/go-mdns/mdns/dns"
	"github.com/spf13/cobra"
)

const (
	typeParamStr = "type"
)

var queryTypes = map[string]mdns.Type{
	"PTR":  mdns.PTR,
	"SRV":  mdns.SRV,
	"TXT":  mdns.TXT,
	"A":    mdns.A,
	"AAAA": mdns.AAAA,
	"ANY":  mdns.ANY,
}

func allSupportedQueryTypes() []string {
	return []string{"PTR", "SRV", "TXT", "A", "AAAA", "ANY"}
}

func newQueryTypeFromString(s string) (mdns.Type, error) {
	t, ok := queryTypes[strings.ToUpper(strings.TrimSpace(s))]
	if !ok {
		return mdns.ANY, fmt.Errorf("invalid query type: %s", s)
	}
	return t, nil
}

var queryCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   "query [name]",
	Short: "Send a single question and print the answering records",
	Long: `Send a single question to the multicast address, and print the records of
every answer until the timeout elapses.

Use it to look at the raw records of a responder. Use "browse" and "resolve" to
discover the services instead.`,
	Example: `  mdnslookup query _matterc._udp.local
  mdnslookup query --type SRV DD200C20D25AE5F7._matterc._udp.local
  mdnslookup query --type ANY --format json macmini.local`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		typeStr, err := cmd.Flags().GetString(typeParamStr)
		if err != nil {
			return err
		}
		queryType, err := newQueryTypeFromString(typeStr)
		if err != nil {
			return err
		}
		unicast, err := cmd.Flags().GetBool(unicastParamStr)
		if err != nil {
			return err
		}

		queryName := mdns.NewQuery(
			mdns.WithQueryService(mdns.DefaultQueryService),
		).Name()
		if 0 < len(args) {
			queryName = args[0]
		}

		client, stop, err := startClient()
		if err != nil {
			return err
		}
		defer stop()

		var outputErr error
		query := mdns.NewQuery(
			mdns.WithQueryName(queryName),
			mdns.WithQueryType(queryType),
			mdns.WithQueryUnicastResponse(unicast),
			mdns.WithQueryMessageHandler(func(msg dns.Message) {
				if !msg.IsResponse() {
					return
				}
				if err := outputMessage(msg); err != nil {
					outputErr = err
				}
			}),
		)

		ctx, cancel := context.WithTimeout(cmd.Context(), queryTimeout())
		defer cancel()

		if _, err := client.Query(ctx, query); err != nil {
			return err
		}

		return outputErr
	},
}

func init() {
	queryCmd.Flags().String(typeParamStr, "PTR", fmt.Sprintf("question record type: %s", strings.Join(allSupportedQueryTypes(), "|")))
	queryCmd.Flags().Bool(unicastParamStr, false, "request unicast responses (QU)")
	rootCmd.AddCommand(queryCmd)
}
