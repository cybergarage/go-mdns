package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cybergarage/go-mdns/mdns"
	"github.com/cybergarage/go-mdns/mdns/dns"
)

func TestWriteRecordTableEscapesFields(t *testing.T) {
	const printable = "caf\u00e9"
	record := dns.NewQuestion(dns.WithQuestionName(printable + "\npart\tname.local"))
	var out bytes.Buffer
	writeRecordTable(&out, mdns.ResourceRecordSet{record})
	if got := out.String(); strings.Count(got, "\n") != 1 || !strings.Contains(got, printable+`\npart\tname.local`) {
		t.Fatalf("unexpected table row: %q", got)
	}
}
