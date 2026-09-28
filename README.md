# go-mdns

![](https://img.shields.io/badge/status-Work%20In%20Progress-8A2BE2)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-mdns)
[![test](https://github.com/cybergarage/go-mdns/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-mdns/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-mdns.svg)](https://pkg.go.dev/github.com/cybergarage/go-mdns)
 [![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/cybergarage/go-mdns)
 [![codecov](https://codecov.io/gh/cybergarage/go-mdns/graph/badge.svg?token=OCU5V0H3OX)](https://codecov.io/gh/cybergarage/go-mdns)

go-mdns is a Go library for Multicast DNS (mDNS) and DNS Service Discovery (DNS-SD) as defined in [RFC 6762](https://www.rfc-editor.org/rfc/rfc6762) and [RFC 6763](https://www.rfc-editor.org/rfc/rfc6763).

**Note:** 🌱 This is a spare-time hobby project, so progress may be slow and changes may appear irregular. Thank you for your patience 🙂

## Status

go-mdns is a **client (querier)** and a **server (responder)** library. The client browses, resolves and looks up the services and the hosts which other responders advertise, and the server publishes the services of this node.

**The responder does not probe for name conflicts yet** (RFC 6762, 8.1 and 9): a registered instance name and host name are assumed to be unique on the link. Probing and conflict resolution are planned for v1.0.0.

| | Status |
| --- | --- |
| Browsing, resolving and host lookup (querier) | Supported |
| Registering, announcing and answering for a service (responder) | Supported, without probing |
| `mdnslookup` command | Supported |
| `mdnsd` command | Supported |

### What the client supports

- Browsing a service type continuously, with the added, updated and removed changes (RFC 6763, 4)
- Selective instance enumeration with the subtypes, such as the Matter `_S<n>` and `_L<n>` subtypes (RFC 6763, 7.1)
- Service type enumeration through `_services._dns-sd._udp` (RFC 6763, 7.2)
- Service instance resolution with the SRV and the TXT records (RFC 6763, 5)
- Host name resolution with the A and the AAAA records (RFC 6762)
- IPv6 scoped addressing, so that the link-local addresses of a discovered service can be used to connect to it (RFC 4007)
- Query retransmission with a doubled interval (RFC 6762, 5.2)
- Goodbye packets and the TTL expiration of the cached services (RFC 6762, 10.1)
- Selecting the network interfaces and the address families to listen on

### What the server supports

- Registering and deregistering a service with its subtypes, host, port and TXT strings (RFC 6763)
- Answering the PTR queries for the service type, its subtypes and the service type enumeration, the SRV and TXT queries for the instance, the A and AAAA queries for the host, and ANY; a PTR answer carries the SRV, TXT and address records as additional records (RFC 6763, 12)
- Publishing the addresses of the interface a query arrives on, unless the service is given its own
- Announcing a service twice when it is registered or the server starts, and sending goodbye records when it is deregistered or the server stops (RFC 6762, 8.3 and 10.1)
- Known-Answer suppression of the answers the querier already holds (RFC 6762, 7.1)
- Unicast responses to QU queries (RFC 6762, 5.4) and to legacy unicast queries from a port other than 5353 (6.7)
- The cache-flush bit on the unique records, and a random 20-120 ms delay of a multicast response with shared records (RFC 6762, 6 and 10.2)

### What is not supported yet

- Probing and conflict resolution of the registered names (RFC 6762, 8.1 and 9)
- Negative responses with NSEC records (RFC 6762, 6.1)
- Name compression when a message is written (the compression pointers are resolved when a message is read)
- Truncated messages (the TC bit) and the Known-Answer list continuation (RFC 6762, 7.2)
- Known-Answer lists in the client's own queries, and duplicate question suppression (RFC 6762, 7.1 - 7.4)
- Following the interface changes, such as a link going up or down, while the client or the server is running

The `mdns/dns` package is exported so that the records of a message can be read, but its API is not stable until v1.0.0.

## Install

```
go get -u github.com/cybergarage/go-mdns
```

The `mdnslookup` command is installed with:

```
go install github.com/cybergarage/go-mdns/cmd/mdnslookup@latest
```

## Usage

### Browsing the services

`Browse` keeps querying, and reports a service when it is added, updated or removed. It blocks until the context is done.

```go
client := mdns.NewClient()
if err := client.Start(); err != nil {
	return err
}
defer client.Stop()

query := mdns.NewQuery(
	mdns.WithQueryService("_matterc._udp"),
)

err := client.Browse(context.Background(), query, func(event mdns.ServiceEvent) {
	service := event.Service
	log.Printf("%s %s (%s:%d)", event.Type, service.FullName(), service.Host(), service.Port())
})
```

Use `mdns.WithQuerySubtype()` to browse a subtype, such as `_S3._sub._matterc._udp`:

```go
query := mdns.NewQuery(
	mdns.WithQuerySubtype("_S3"),
	mdns.WithQueryService("_matterc._udp"),
)
```

### Resolving a service instance

`Resolve` returns the host, the port, the addresses and the TXT attributes of a service instance name.

```go
service, err := client.Resolve(context.Background(), "DD200C20D25AE5F7._matterc._udp.local")
if err != nil {
	return err
}

for _, addr := range service.Addrs() {
	// addr holds the service port, and the IPv6 scoped addressing zone
	// is set for the link-local addresses, such as [fe80::1%en0]:5540.
	log.Println(addr.String())
}

if attr, ok := service.LookupResourceAttribute("CM"); ok {
	log.Println(attr.Value())
}
```

### Resolving a host name

```go
addrs, err := client.LookupHost(context.Background(), "macmini.local")
```

### Publishing a service

`Server.Register` publishes a service: the server announces it, answers the queries for it, and withdraws it with goodbye records when it is deregistered or the server stops.

```go
server := mdns.NewServer()
if err := server.Start(); err != nil {
	return err
}
defer server.Stop()

err := server.Register(&mdns.LocalService{
	Instance: "665F6E75B5D3A9C2",
	Service:  "_matterc._udp",
	Subtypes: []string{"_L3840", "_S15", "_V65521", "_CM"},
	Host:     "B75AFB458ECD6D6F",
	Port:     5540,
	TXT:      []string{"D=3840", "CM=1", "VP=65521+32769"},
})
```

Without `Addresses`, the host name resolves to the addresses of the interface a query arrives on. Register the service again to update it, such as its TXT strings.

### Selecting the interfaces

```go
ifi, _ := net.InterfaceByName("en0")

client := mdns.NewClient(
	mdns.WithClientInterfaces(ifi),
	mdns.WithClientIPv4Enabled(false),
	mdns.WithClientQueryTimeout(10*time.Second),
)
```

## Command

`mdnslookup` browses and resolves the services from a terminal.

```
$ mdnslookup types
$ mdnslookup browse _matterc._udp
$ mdnslookup browse --subtype _S3 --duration 10s _matterc._udp
$ mdnslookup resolve DD200C20D25AE5F7._matterc._udp.local
$ mdnslookup host macmini.local
$ mdnslookup query --type SRV DD200C20D25AE5F7._matterc._udp.local
$ mdnslookup monitor
$ mdnslookup decode mdnstest/dumps/matter-answer-01.dump
```

Every command takes `--format table|json|csv`, and the common `--interface`, `--family` and `--timeout` flags. `decode` needs no network, so a message which was captured elsewhere can be analyzed offline.

`mdnsd` publishes a service from a terminal until it is interrupted.

```
$ mdnsd -name demo -service _http._tcp -port 8080 -txt path=/
```

## Testing

`make test` runs the unit tests and the tests on the network. When the mDNS clients of the operating system are installed, the tests in `mdnstest/interop_test.go` also publish a service with the responder and check that the clients find it:

| Client | Checked |
| --- | --- |
| `dns-sd` (Bonjour: macOS, Windows) | `-B` by type and by subtype, `-L`, `-G`, and the removal after a goodbye |
| `avahi-browse`, `avahi-resolve` (Avahi: Linux) | `--resolve` by type and by subtype, `--name`, and the removal after a goodbye |

A test is skipped when its client is not installed or, for Avahi, when `avahi-daemon` is not running. `go test -short` skips them all. `TestInteropClients` logs which clients were found, so the verbose log shows which of the tests ran:

```
$ make test-interop
    interop_test.go:177: dns-sd         /usr/bin/dns-sd: Currently running daemon (system service) is version mDNSResponder-...
    interop_test.go:177: avahi-browse   not installed: its tests are skipped
```

Set `GO_MDNS_TEST_REQUIRE` to the clients which must be available, such as `avahi` or `avahi,dns-sd`, to fail instead of skipping the tests of a missing one. The GitHub Actions workflow installs Avahi on Ubuntu and sets `GO_MDNS_TEST_REQUIRE=avahi`. The Homebrew `avahi` formula is Linux only, so on macOS the tests use `dns-sd`.

`go test` caches the results, and a cached result does not run the clients again; use `make test`, `make test-interop` or `go test -count=1`.

# User Guides

- Operation
  - [mdnslookup](doc/mdnslookup.md)
  - [mdnsd](doc/mdnsd.md)

## References

### DNS
- [RFC 1034: DOMAIN NAMES - CONCEPTS AND FACILITIES](https://www.rfc-editor.org/rfc/rfc1034)
- [RFC 1035: DOMAIN NAMES - IMPLEMENTATION AND SPECIFICATION](https://www.rfc-editor.org/rfc/rfc1035)
- [RFC 2782: A DNS RR for specifying the location of services (DNS SRV)](https://www.rfc-editor.org/rfc/rfc2782)
- [RFC 3901: DNS IPv6 Transport Operational Guidelines](https://www.rfc-editor.org/rfc/rfc3901)
- [RFC 4034: Resource Records for the DNS Security Extensions](https://datatracker.ietf.org/doc/html/rfc4034)

### mDNS and DNS-SD

- [RFC 6762: Multicast DNS](https://www.rfc-editor.org/rfc/rfc6762)
- [RFC 6763: DNS-Based Service Discovery](https://www.rfc-editor.org/rfc/rfc6763)
