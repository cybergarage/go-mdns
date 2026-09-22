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
	"strings"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mdns/mdns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ProgramName = "mdnslookup"

	VerboseParamStr   = "verbose"
	DebugParamStr     = "debug"
	InterfaceParamStr = "interface"
	FamilyParamStr    = "family"
	TimeoutParamStr   = "timeout"
)

const (
	FamilyAllStr  = "all"
	FamilyIPv4Str = "ipv4"
	FamilyIPv6Str = "ipv6"
)

func allSupportedFamilies() []string {
	return []string{
		FamilyAllStr,
		FamilyIPv4Str,
		FamilyIPv6Str,
	}
}

var rootCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   ProgramName,
	Short: "Browse and resolve mDNS (DNS-SD) services on the local link",
	Long: `mdnslookup browses and resolves the Multicast DNS (RFC 6762) and
DNS-Based Service Discovery (RFC 6763) services on the local link.

  mdnslookup types                           list the advertised service types
  mdnslookup browse _matterc._udp            browse the instances of a service type
  mdnslookup resolve <instance>._matterc._udp.local
                                             resolve an instance to a host and a port
  mdnslookup host <device>.local             resolve a host name to its addresses
  mdnslookup monitor                         watch the mDNS messages on the link
  mdnslookup decode <file>                   decode a recorded mDNS message

This tool is a client. Advertising a service is not supported yet.`,
	Version:           mdns.Version,
	DisableAutoGenTag: true,
	// The usage is not printed for a runtime error, such as a name which is
	// not found, because the command line itself is valid.
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		log.SetDefault(nil)
		verbose := viper.GetBool(VerboseParamStr)
		debug := viper.GetBool(DebugParamStr)
		if debug {
			verbose = true
		}
		if verbose {
			if debug {
				log.SetDefault(log.NewStdoutLogger(log.LevelDebug))
			} else {
				log.SetDefault(log.NewStdoutLogger(log.LevelInfo))
			}
			log.Infof("%s version %s", ProgramName, mdns.Version)
			log.Infof("verbose:%t, debug:%t", verbose, debug)
		}
		return nil
	},
}

func GetRootCommand() *cobra.Command {
	return rootCmd
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	viper.SetEnvPrefix(ProgramName)

	viper.SetDefault(FormatParamStr, FormatTableStr)
	rootCmd.PersistentFlags().String(FormatParamStr, FormatTableStr, fmt.Sprintf("output format: %s", strings.Join(allSupportedFormats(), "|")))
	viper.BindPFlag(FormatParamStr, rootCmd.PersistentFlags().Lookup(FormatParamStr))
	viper.BindEnv(FormatParamStr)

	viper.SetDefault(InterfaceParamStr, []string{})
	rootCmd.PersistentFlags().StringSliceP(InterfaceParamStr, "i", []string{}, "network interfaces to use (all available interfaces by default)")
	viper.BindPFlag(InterfaceParamStr, rootCmd.PersistentFlags().Lookup(InterfaceParamStr))
	viper.BindEnv(InterfaceParamStr)

	viper.SetDefault(FamilyParamStr, FamilyAllStr)
	rootCmd.PersistentFlags().String(FamilyParamStr, FamilyAllStr, fmt.Sprintf("address family to use: %s", strings.Join(allSupportedFamilies(), "|")))
	viper.BindPFlag(FamilyParamStr, rootCmd.PersistentFlags().Lookup(FamilyParamStr))
	viper.BindEnv(FamilyParamStr)

	viper.SetDefault(TimeoutParamStr, mdns.DefaultQueryTimeout)
	rootCmd.PersistentFlags().DurationP(TimeoutParamStr, "t", mdns.DefaultQueryTimeout, "query timeout")
	viper.BindPFlag(TimeoutParamStr, rootCmd.PersistentFlags().Lookup(TimeoutParamStr))
	viper.BindEnv(TimeoutParamStr)

	viper.SetDefault(VerboseParamStr, false)
	rootCmd.PersistentFlags().Bool(VerboseParamStr, false, "enable verbose output")
	viper.BindPFlag(VerboseParamStr, rootCmd.PersistentFlags().Lookup(VerboseParamStr))
	viper.BindEnv(VerboseParamStr)

	viper.SetDefault(DebugParamStr, false)
	rootCmd.PersistentFlags().Bool(DebugParamStr, false, "enable debug output")
	viper.BindPFlag(DebugParamStr, rootCmd.PersistentFlags().Lookup(DebugParamStr))
	viper.BindEnv(DebugParamStr)
}

// queryTimeout returns the configured query timeout.
func queryTimeout() time.Duration {
	timeout := viper.GetDuration(TimeoutParamStr)
	if timeout <= 0 {
		return mdns.DefaultQueryTimeout
	}
	return timeout
}
