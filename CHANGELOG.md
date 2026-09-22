# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.9.0]

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
