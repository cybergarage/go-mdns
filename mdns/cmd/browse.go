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
	"sync"
	"syscall"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns"
	"github.com/spf13/cobra"
)

const (
	subtypeParamStr  = "subtype"
	domainParamStr   = "domain"
	unicastParamStr  = "unicast"
	durationParamStr = "duration"
	resolveParamStr  = "resolve"
)

var browseCmd = &cobra.Command{ // nolint:exhaustruct,exhaustruct_v5
	Use:   "browse [service]",
	Short: "Browse the instances of a service type",
	Long: `Browse the instances of a service type, and report the instances as they are
added, updated and removed. The browse runs until it is interrupted, or until
the duration elapses when --duration is set. Without a service type, it browses
the service types which are advertised on the link.

A responder usually answers a browse with the host, the port and the addresses
of each instance, but one may answer only with the instance names. With
--resolve, such an instance is resolved before it is reported, as
avahi-browse --resolve does.`,
	Example: `  mdnslookup browse _matterc._udp
  mdnslookup browse --subtype _S3 _matterc._udp
  mdnslookup browse --resolve _http._tcp
  mdnslookup browse --duration 10s --format json _matter._tcp
  mdnslookup browse`,
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

		resolve, err := cmd.Flags().GetBool(resolveParamStr)
		if err != nil {
			return err
		}
		// The service types have nothing to resolve.
		resolve = resolve && service != ""

		var resolving sync.WaitGroup
		defer resolving.Wait()

		return client.Browse(ctx, query, func(event mdns.ServiceEvent) {
			if !resolve || event.Type == mdns.ServiceRemoved || isResolved(event.Service) {
				writer.Write(event.Type.String(), event.Service)
				return
			}
			// The browse handler must not block the messages which the
			// resolution waits for.
			resolving.Go(func() {
				writer.Write(event.Type.String(), resolveService(ctx, client, event.Service))
			})
		})
	},
}

// isResolved reports whether service holds what a resolution finds: its host,
// port and addresses.
func isResolved(service mdns.Service) bool {
	return service.Host() != "" && 0 < service.Port() && 0 < len(service.Addresses())
}

// resolveService resolves service by its instance name, and returns it as it
// is when the resolution fails.
func resolveService(ctx context.Context, client mdns.Client, service mdns.Service) mdns.Service {
	resolveCtx, cancel := context.WithTimeout(ctx, queryTimeout())
	defer cancel()
	resolved, err := client.Resolve(resolveCtx, service.FullName())
	if err != nil {
		log.Warnf("resolve %s: %s", service.FullName(), err)
		return service
	}
	return resolved
}

func init() {
	browseCmd.Flags().String(subtypeParamStr, "", "service subtype to browse (_S3, _L840, ...)")
	browseCmd.Flags().String(domainParamStr, mdns.DefaultQueryDomain, "domain to browse")
	browseCmd.Flags().Bool(unicastParamStr, false, "request unicast responses (QU)")
	browseCmd.Flags().Duration(durationParamStr, time.Duration(0), "browse duration (0 means until interrupted)")
	browseCmd.Flags().Bool(resolveParamStr, false, "resolve an instance which is reported without its host, port or addresses")
	rootCmd.AddCommand(browseCmd)
}
