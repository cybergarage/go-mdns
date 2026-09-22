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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/cybergarage/go-mdns/mdns"
	"github.com/spf13/viper"
)

// outputFormat returns the configured output format.
func outputFormat() (Format, error) {
	return NewFormatFromString(viper.GetString(FormatParamStr))
}

// serviceObject represents a service for the JSON and the CSV output.
type serviceObject struct {
	Event      string   `json:"event,omitempty"`
	Name       string   `json:"name"`
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	Interface  string   `json:"interface,omitempty"`
	Addresses  []string `json:"addresses"`
	Attributes []string `json:"attributes"`
}

func addrString(addr *net.UDPAddr) string {
	host := addr.IP.String()
	if 0 < len(addr.Zone) {
		host += "%" + addr.Zone
	}
	return host
}

func newServiceObject(service mdns.Service) serviceObject {
	addrs := make([]string, 0)
	for _, addr := range service.Addrs() {
		addrs = append(addrs, addrString(addr))
	}

	attrs := make([]string, 0)
	for _, attr := range service.ResourceAttributes() {
		attrs = append(attrs, attr.String())
	}

	ifname := ""
	if ifi := service.Interface(); ifi != nil {
		ifname = ifi.Name
	}

	return serviceObject{
		Event:      "",
		Name:       service.FullName(),
		Host:       service.Host(),
		Port:       service.Port(),
		Interface:  ifname,
		Addresses:  addrs,
		Attributes: attrs,
	}
}

func (obj serviceObject) record(withEvent bool) []string {
	record := []string{}
	if withEvent {
		record = append(record, obj.Event)
	}
	return append(record,
		obj.Name,
		obj.Host,
		strconv.Itoa(obj.Port),
		obj.Interface,
		strings.Join(obj.Addresses, " "),
		strings.Join(obj.Attributes, " "),
	)
}

func serviceHeader(withEvent bool) []string {
	header := []string{}
	if withEvent {
		header = append(header, "EVENT")
	}
	return append(header, "NAME", "HOST", "PORT", "INTERFACE", "ADDRESSES", "ATTRIBUTES")
}

// serviceWriter writes the services in the configured output format. The
// services are written as they are found, so that a browse can report a change
// while it is running.
type serviceWriter struct {
	mutex       sync.Mutex
	format      Format
	withEvent   bool
	wroteHeader bool
	tabWriter   *tabwriter.Writer
	csvWriter   *csv.Writer
}

func newServiceWriter(withEvent bool) (*serviceWriter, error) {
	format, err := outputFormat()
	if err != nil {
		return nil, err
	}

	writer := &serviceWriter{
		mutex:       sync.Mutex{},
		format:      format,
		withEvent:   withEvent,
		wroteHeader: false,
		tabWriter:   nil,
		csvWriter:   nil,
	}

	switch format {
	case FormatTable:
		writer.tabWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	case FormatCSV:
		writer.csvWriter = csv.NewWriter(os.Stdout)
	case FormatJSON:
	}

	return writer, nil
}

// Write writes the specified service.
func (writer *serviceWriter) Write(event string, service mdns.Service) error {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	obj := newServiceObject(service)
	obj.Event = event

	switch writer.format {
	case FormatJSON:
		// The services are written as JSON Lines, because they are
		// written as they are found.
		b, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(b))
	case FormatCSV:
		if !writer.wroteHeader {
			if err := writer.csvWriter.Write(serviceHeader(writer.withEvent)); err != nil {
				return err
			}
			writer.wroteHeader = true
		}
		if err := writer.csvWriter.Write(obj.record(writer.withEvent)); err != nil {
			return err
		}
		writer.csvWriter.Flush()
	case FormatTable:
		if !writer.wroteHeader {
			fmt.Fprintln(writer.tabWriter, strings.Join(serviceHeader(writer.withEvent), "\t"))
			writer.wroteHeader = true
		}
		fmt.Fprintln(writer.tabWriter, strings.Join(obj.record(writer.withEvent), "\t"))
		writer.tabWriter.Flush()
	}

	return nil
}

// Flush flushes the buffered output.
func (writer *serviceWriter) Flush() {
	writer.mutex.Lock()
	defer writer.mutex.Unlock()

	if writer.tabWriter != nil {
		writer.tabWriter.Flush()
	}
	if writer.csvWriter != nil {
		writer.csvWriter.Flush()
	}
}

// outputServices writes the specified services in the configured format.
func outputServices(services []mdns.Service) error {
	writer, err := newServiceWriter(false)
	if err != nil {
		return err
	}
	defer writer.Flush()

	for _, service := range services {
		if err := writer.Write("", service); err != nil {
			return err
		}
	}

	return nil
}

// addrObject represents a resolved host address for the JSON output.
type addrObject struct {
	Host      string   `json:"host"`
	Addresses []string `json:"addresses"`
}

// outputHostAddrs writes the specified host addresses in the configured format.
func outputHostAddrs(host string, addrs []*net.IPAddr) error {
	format, err := outputFormat()
	if err != nil {
		return err
	}

	addrStrs := make([]string, 0, len(addrs))
	for _, addr := range addrs {
		addrStrs = append(addrStrs, addr.String())
	}

	switch format {
	case FormatJSON:
		b, err := json.Marshal(addrObject{
			Host:      host,
			Addresses: addrStrs,
		})
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(b))
	case FormatCSV:
		csvWriter := csv.NewWriter(os.Stdout)
		defer csvWriter.Flush()
		if err := csvWriter.Write([]string{"HOST", "ADDRESS"}); err != nil {
			return err
		}
		for _, addr := range addrStrs {
			if err := csvWriter.Write([]string{host, addr}); err != nil {
				return err
			}
		}
	case FormatTable:
		tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer tabWriter.Flush()
		fmt.Fprintln(tabWriter, "HOST\tADDRESS")
		for _, addr := range addrStrs {
			fmt.Fprintf(tabWriter, "%s\t%s\n", host, addr)
		}
	}

	return nil
}

// messageObject represents a DNS message for the JSON output.
type messageObject struct {
	From      string         `json:"from,omitempty"`
	Response  bool           `json:"response"`
	Questions []recordObject `json:"questions"`
	Answers   []recordObject `json:"answers"`
	Authority []recordObject `json:"authority"`
	Additions []recordObject `json:"additions"`
}

// recordObject represents a DNS record for the JSON output.
type recordObject struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	TTL     uint   `json:"ttl"`
	Content string `json:"content"`
}

func newRecordObjects(records mdns.ResourceRecordSet) []recordObject {
	objs := make([]recordObject, 0, len(records))
	for _, record := range records {
		objs = append(objs, recordObject{
			Name:    record.Name(),
			Type:    record.Type().String(),
			TTL:     record.TTL(),
			Content: record.Content(),
		})
	}
	return objs
}

// outputMessage writes the specified message in the configured format.
func outputMessage(msg mdns.Message) error {
	format, err := outputFormat()
	if err != nil {
		return err
	}

	from := ""
	if msg.From() != nil {
		from = msg.From().String()
	}

	questions := mdns.ResourceRecordSet{}
	for _, q := range msg.Questions() {
		questions = append(questions, q)
	}

	switch format {
	case FormatJSON:
		b, err := json.Marshal(messageObject{
			From:      from,
			Response:  msg.IsResponse(),
			Questions: newRecordObjects(questions),
			Answers:   newRecordObjects(msg.Answers()),
			Authority: newRecordObjects(msg.NameServers()),
			Additions: newRecordObjects(msg.Additions()),
		})
		if err != nil {
			return err
		}
		fmt.Fprintln(os.Stdout, string(b))
	case FormatCSV:
		csvWriter := csv.NewWriter(os.Stdout)
		defer csvWriter.Flush()
		if err := csvWriter.Write([]string{"SECTION", "NAME", "TYPE", "TTL", "CONTENT"}); err != nil {
			return err
		}
		sections := []struct {
			name    string
			records mdns.ResourceRecordSet
		}{
			{"question", questions},
			{"answer", msg.Answers()},
			{"authority", msg.NameServers()},
			{"additional", msg.Additions()},
		}
		for _, section := range sections {
			for _, obj := range newRecordObjects(section.records) {
				if err := csvWriter.Write([]string{section.name, obj.Name, obj.Type, strconv.Itoa(int(obj.TTL)), obj.Content}); err != nil {
					return err
				}
			}
		}
	case FormatTable:
		kind := "QUERY"
		if msg.IsResponse() {
			kind = "RESPONSE"
		}
		if 0 < len(from) {
			fmt.Fprintf(os.Stdout, "%s from %s\n", kind, from)
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", kind)
		}
		if 0 < len(questions) {
			fmt.Fprintf(os.Stdout, "%s\n", questions.String())
		}
		records := msg.ResourceRecordSet()
		if 0 < len(records) {
			fmt.Fprintf(os.Stdout, "%s\n", records.String())
		}
		fmt.Fprintln(os.Stdout, "")
	}

	return nil
}
