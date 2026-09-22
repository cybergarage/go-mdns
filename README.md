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

go-mdns is a **client (querier)** library. It browses, resolves and looks up the services and the hosts which other responders advertise.

**The server (responder) side is under development.** `mdns.Server` and the `mdnsd` command are included, but they only listen for the messages: they do not register a service, answer a query or announce anything yet. Do not use them to advertise a service. The responder is planned for v1.0.0.

| | Status |
| --- | --- |
| Browsing, resolving and host lookup (querier) | Supported |
| Registering and announcing a service (responder) | **Under development** |
| `mdnslookup` command | Supported |
| `mdnsd` command | **Under development** |

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

### What is not supported yet

- Registering, probing, announcing and answering a query (the responder side)
- Name compression when a message is written (the compression pointers are resolved when a message is read)
- Truncated messages (the TC bit) and the Known-Answer list continuation (RFC 6762, 7.2)
- Known-Answer suppression and duplicate question suppression (RFC 6762, 7.1 - 7.4)
- Legacy unicast queries from a source port other than 5353 (RFC 6762, 6.7)
- Following the interface changes, such as a link going up or down, while the client is running

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

# User Guides

- Operation
  - [mdnslookup](doc/mdnslookup.md)
  - [mdnsd](doc/mdnsd.md) (under development)

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
