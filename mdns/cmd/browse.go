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
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cybergarage/go-mdns/mdns"
	"github.com/spf13/cobra"
)

const (
	subtypeParamStr  = "subtype"
	domainParamStr   = "domain"
	unicastParamStr  = "unicast"
	durationParamStr = "duration"
)

var browseCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "browse [service]",
	Short: "Browse the instances of a service type",
	Long: `Browse the instances of a service type, and report the instances as they are
added, updated and removed. The browse runs until it is interrupted, or until
the duration elapses when --duration is set.`,
	Example: `  mdnslookup browse _matterc._udp
  mdnslookup browse --subtype _S3 _matterc._udp
  mdnslookup browse --duration 10s --format json _matter._tcp`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := ""
		if 0 < len(args) {
			service = args[0]
		}

		subtype, err := cmd.Flags().GetString(subtypeParamStr)
		if err != nil {
			return err
		}
		domain, err := cmd.Flags().GetString(domainParamStr)
		if err != nil {
			return err
		}
		unicast, err := cmd.Flags().GetBool(unicastParamStr)
		if err != nil {
			return err
		}
		duration, err := cmd.Flags().GetDuration(durationParamStr)
		if err != nil {
			return err
		}

		client, stop, err := startClient()
		if err != nil {
			return err
		}
		defer stop()

		writer, err := newServiceWriter(true)
		if err != nil {
			return err
		}
		defer writer.Flush()

		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		if 0 < duration {
			var timeoutCancel context.CancelFunc
			ctx, timeoutCancel = context.WithTimeout(ctx, duration)
			defer timeoutCancel()
		}

		query := mdns.NewQuery(
			mdns.WithQuerySubtype(subtype),
			mdns.WithQueryService(service),
			mdns.WithQueryDomain(domain),
			mdns.WithQueryUnicastResponse(unicast),
		)

		return client.Browse(ctx, query, func(event mdns.ServiceEvent) {
			writer.Write(event.Type.String(), event.Service)
		})
	},
}

func init() {
	browseCmd.Flags().String(subtypeParamStr, "", "service subtype to browse (_S3, _L840, ...)")
	browseCmd.Flags().String(domainParamStr, mdns.DefaultQueryDomain, "domain to browse")
	browseCmd.Flags().Bool(unicastParamStr, false, "request unicast responses (QU)")
	browseCmd.Flags().Duration(durationParamStr, time.Duration(0), "browse duration (0 means until interrupted)")
	rootCmd.AddCommand(browseCmd)
}
