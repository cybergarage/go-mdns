# Using go-mdns

go-mdns is a client (querier) library for Multicast DNS ([RFC 6762](https://www.rfc-editor.org/rfc/rfc6762)) and DNS-Based Service Discovery ([RFC 6763](https://www.rfc-editor.org/rfc/rfc6763)). This guide covers the four things a client does: enumerating the service types, browsing a service type, resolving an instance, and resolving a host name.

The responder side is under development, so go-mdns cannot advertise a service yet. See the [README](../README.md#status) for the current status.

## The client

A client listens on the mDNS multicast address of every available network interface, and it must be started before a query is sent.

```go
client := mdns.NewClient()
if err := client.Start(); err != nil {
	return err
}
defer client.Stop()
```

`NewClient` takes options to narrow what the client listens on, and to change the timing of the queries.

```go
client := mdns.NewClient(
	mdns.WithClientInterfaces(ifi),        // listen only on this interface
	mdns.WithClientIPv4Enabled(false),     // listen only on the IPv6 addresses
	mdns.WithClientQueryTimeout(10*time.Second),
	mdns.WithClientQueryInterval(time.Second),      // RFC 6762, 5.2
	mdns.WithClientMaxQueryInterval(time.Minute),
)
```

Selecting an interface matters when a link-local address is involved: an address which is discovered on one link cannot be reached on another.

## Building a query

A query names what to ask for. The name is built from the subtype, the service and the domain, or it is given as it is.

```go
// _matterc._udp.local
mdns.NewQuery(
	mdns.WithQueryService("_matterc._udp"),
)

// _S3._sub._matterc._udp.local (RFC 6763, 7.1)
mdns.NewQuery(
	mdns.WithQuerySubtype("_S3"),
	mdns.WithQueryService("_matterc._udp"),
)

// A name as it is
mdns.NewQuery(
	mdns.WithQueryName("DD200C20D25AE5F7._matterc._udp.local"),
	mdns.WithQueryType(mdns.SRV),
)
```

The question type is `PTR` by default, which is what RFC 6763 browses a service type with. `WithQueryType` selects another type, and `WithQueryUnicastResponse(true)` sets the unicast response bit (QU) of RFC 6762 (5.4). The bit is not set by default, because a multicast response lets the other nodes on the link update their caches.

## Enumerating the service types

RFC 6763 (7.2) enumerates the advertised service types by browsing `_services._dns-sd._udp`.

```go
query := mdns.NewQuery(
	mdns.WithQueryService(mdns.ServiceTypeEnumerationName),
)

services, err := client.Query(context.Background(), query)
for _, service := range services {
	fmt.Println(service.FullName()) // _matterc._udp.local
}
```

`Query` sends a one-shot query, and returns the services which answered before the context is done. When the context has no deadline, the client query timeout is used.

## Browsing a service type

A one-shot query is a snapshot. `Browse` keeps querying, and reports a service when it is added, updated or removed, which is what a controller needs to follow the nodes on a link.

```go
query := mdns.NewQuery(
	mdns.WithQueryService("_matter._tcp"),
)

err := client.Browse(ctx, query, func(event mdns.ServiceEvent) {
	switch event.Type {
	case mdns.ServiceAdded:
	case mdns.ServiceUpdated:
	case mdns.ServiceRemoved:
	}
	fmt.Println(event.Service.FullName())
})
```

`Browse` blocks until the context is done. A service is removed when its responder announces a goodbye packet with a zero TTL (RFC 6762, 10.1), or when the shortest TTL of its records elapses.

The query is retransmitted while browsing, and the interval is doubled up to the maximum interval (RFC 6762, 5.2), so a device which is not ready to answer the first query is still found.

## Resolving a service instance

A browse gives a service instance name. `Resolve` turns it into something to connect to.

```go
service, err := client.Resolve(ctx, "DD200C20D25AE5F7._matterc._udp.local")
if err != nil {
	return err
}

service.Host()  // 84FCE6036F38.local
service.Port()  // 5540
```

`Resolve` queries the SRV and the TXT records of the name, and it queries the address records of the SRV target when the responder does not send them with the answer (RFC 6763, 5). It returns as soon as it is answered, so it does not wait for the whole timeout.

### Addresses

`Addresses` returns the bare IP addresses, and `Addrs` returns them with the service port and the IPv6 scoped addressing zone.

```go
for _, addr := range service.Addrs() {
	// [fe80::86fc:e6ff:fe03:6f38%en0]:5540
	conn, err := net.DialUDP("udp", nil, addr)
}
```

An IPv6 link-local address is ambiguous without its zone ([RFC 4007](https://www.rfc-editor.org/rfc/rfc4007)). go-mdns sets the zone from the interface which the response was received on, which `service.Interface()` also returns. Use `Addrs` whenever a link-local address may be involved, which is the common case for the devices which advertise themselves on a link.

### TXT attributes

The TXT record attributes follow RFC 6763 (6). The keys are case insensitive, the first occurrence of a key wins, and an attribute without `=` is present with no value.

```go
if attr, ok := service.LookupResourceAttribute("CM"); ok {
	attr.Value()     // "2"
	attr.HasValue()  // true
}

for _, attr := range service.ResourceAttributes() {
	fmt.Println(attr.String())
}
```

### Records

The raw records of the response are available when a detail is needed which the service does not expose.

```go
for _, record := range service.ResourceRecordSet() {
	fmt.Println(record.Name(), record.Type(), record.TTL(), record.Content())
}
```

## Resolving a host name

Resolving a host name is the core of Multicast DNS, and it needs no service records.

```go
addrs, err := client.LookupHost(ctx, "macmini.local")
for _, addr := range addrs {
	fmt.Println(addr.String()) // fe80::46d:889b:988:3dfc%en0
}
```

## Reading the messages

A handler receives every message which the client reads, including the queries and the responses of the other nodes on the link. Use it to watch the link, or to look at a record which the service API does not expose.

```go
client.RegisterMessageHandler(func(msg mdns.Message) {
	if !msg.IsResponse() {
		return
	}
	fmt.Println(msg.String())
})
```

A handler is called from the listener goroutines, so it must not block, and it must guard whatever it shares with the rest of the program.

## The command

The `mdnslookup` command does the same from a terminal, and it is the quickest way to see what a link advertises. See the [command reference](mdnslookup.md).

```
$ mdnslookup types
$ mdnslookup browse _matterc._udp
$ mdnslookup resolve DD200C20D25AE5F7._matterc._udp.local
$ mdnslookup host macmini.local
$ mdnslookup monitor --responses
$ mdnslookup decode mdnstest/dumps/matter-answer-01.dump
```
