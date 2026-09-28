# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

The responder. `Server` publishes services, which makes go-mdns usable to advertise, such as a Matter device advertising itself for commissioning. Probing and conflict resolution remain for v1.0.0.

### Added

- `Server.Register()` and `Server.Deregister()` publish and withdraw a `LocalService`: an instance, a service type, its subtypes, a host, a port, TXT strings and optionally its own addresses. `Server.LocalServices()` lists them.
- The server answers the PTR queries for a registered service type, its subtypes and the service type enumeration, the SRV and TXT queries for the instance, the A and AAAA queries for the host, and ANY. A PTR answer carries the SRV, TXT and address records as additional records (RFC 6763, 12).
- A registered service is announced twice when it is registered or the server starts, and goodbye records are sent when it is deregistered or the server stops (RFC 6762, 8.3 and 10.1).
- Known-Answer suppression, the cache-flush bit on the unique records, and a random 20-120 ms delay of a multicast response with shared records (RFC 6762, 7.1, 10.2 and 6).
- A QU query (RFC 6762, 5.4) and a legacy unicast query from a port other than 5353 (6.7) are answered by unicast; the reply to a legacy query echoes its ID and question, with TTLs of at most 10 seconds.
- `dns.NewPTRResourceRecord()`, `dns.NewSRVResourceRecord()`, `dns.NewTXTResourceRecord()`, `dns.NewAResourceRecord()`, `dns.NewAAAAResourceRecord()` and `dns.NewAddressResourceRecord()` build records, and `dns.WithMessageID()`, `dns.WithMessageAnswers()`, `dns.WithMessageNameServers()` and `dns.WithMessageAdditions()` build messages. `dns.CacheFlush` names the cache-flush bit.
- `mdnsd` publishes the service given by `-name`, `-service`, `-port`, `-host`, `-subtype` and `-txt`.

## [0.9.1] - 2026-09-23

### Fixed

- A service lost its addresses when the response which carried it abbreviated the names of its records. The address records are selected by the SRV target name, so that a response which holds more than one instance does not report the addresses of the other hosts, and a responder which answers with records whose name does not match the target, such as the example of the Matter specification (1.2, 4.3.1.13), left the service without an address at all. The records of the target are still preferred, and the rest are collected only when none of them names the target.
- `make version` overwrote `version.go` with the next patch of the latest tag, so a minor release which was prepared by hand before its tag was made was lost on the next build, and a repository with no tag at all got a version of `..1`. `version.gen` now fails instead of writing a version it cannot derive, and the version is kept when the generated one is not newer.

## [0.9.0] - 2026-09-22

The first release with a client API which covers browsing, resolving and host lookup. The responder side is still under development, and it is planned for v1.0.0.

### Added

- `Client.Browse()` browses a service type continuously, and reports a service when it is added, updated or removed. A service is removed by a goodbye packet with a zero TTL (RFC 6762, 10.1) or by the expiration of its records.
- `Client.Resolve()` resolves a service instance name to its host, port, addresses and TXT attributes (RFC 6763, 5).
- `Client.LookupHost()` resolves a host name to its addresses.
- `Service.Addrs()` returns the service addresses with the service port, and the IPv6 scoped addressing zone is set for the link-local addresses (RFC 4007). `Service.Interface()` returns the interface which the service was discovered on.
- `Service.FullName()` and `Service.TTL()`.
- `WithQueryType()` selects the question record type, and `WithQueryName()` asks for a name as it is.
- `WithQueryUnicastResponse()` sets the unicast response bit (QU).
- `NewClient()` takes options to select the network interfaces and the address families to listen on, and to set the query timeout and the retransmission intervals.
- A query is retransmitted with a doubled interval while it waits for the responses (RFC 6762, 5.2).
- `Attribute.HasValue()` reports whether a TXT attribute has a value.
- `mdnslookup` gained the `types`, `browse`, `resolve`, `host`, `monitor` and `decode` commands, and every command supports `--format table|json|csv` and the common `--interface`, `--family` and `--timeout` flags.
- `FuzzParseMessage` covers the message parser against the malformed messages.

### Changed

- The default question type of a query is `PTR` instead of `ANY`, which is what RFC 6763 browses a service type with.
- The unicast response bit (QU) is no longer set by default, because a multicast response lets the other nodes on the link update their caches (RFC 6762, 5.4).
- `mdnslookup scan` and `mdnslookup dump` are merged into `mdnslookup monitor`.

### Fixed

- The SRV target was read as a single label, so a host was reported as `host` instead of `host.local`, and a compressed target was truncated.
- Every A and AAAA record of a response was collected as an address of the parsed service, so a response which held more than one instance reported the addresses of the other hosts. The address records are now selected by the SRV target name.
- The TXT attributes were parsed by splitting a string at every `=`, so a value which held `=` was rejected, and an attribute without `=` was dropped (RFC 6763, 6.3). The attribute keys are looked up case insensitively, and the first occurrence wins (6.4).
- A name reader followed the compression pointers without a limit, so a message which held pointers referring to each other made the parser loop forever.
- The UDP connection was closed while the listener goroutine was reading it, and the message handlers were called while the handler list was being updated, which `go test -race` reported as data races.
