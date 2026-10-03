## mdnslookup

Browse and resolve mDNS (DNS-SD) services on the local link

### Synopsis

mdnslookup browses and resolves the Multicast DNS (RFC 6762) and
DNS-Based Service Discovery (RFC 6763) services on the local link.

  mdnslookup types                           list the advertised service types
  mdnslookup browse _matterc._udp            browse the instances of a service type
  mdnslookup resolve <instance>._matterc._udp.local
                                             resolve an instance to a host and a port
  mdnslookup host <device>.local             resolve a host name to its addresses
  mdnslookup monitor                         watch the mDNS messages on the link
  mdnslookup decode <file>                   decode a recorded mDNS message

This tool is a client. Use mdnsd to advertise a service.

### Options

```
      --debug               enable debug output
      --family string       address family to use: all|ipv4|ipv6 (default "all")
      --format string       output format: table|json|csv (default "table")
  -h, --help                help for mdnslookup
  -i, --interface strings   network interfaces to use (all available interfaces by default)
  -t, --timeout duration    query timeout (default 5s)
      --verbose             enable verbose output
```

### SEE ALSO

* [mdnslookup browse](#mdnslookup-browse)	 - Browse the instances of a service type
* [mdnslookup decode](#mdnslookup-decode)	 - Decode a recorded mDNS message
* [mdnslookup host](#mdnslookup-host)	 - Resolve a host name to its addresses
* [mdnslookup monitor](#mdnslookup-monitor)	 - Watch the mDNS messages on the link
* [mdnslookup query](#mdnslookup-query)	 - Send a single question and print the answering records
* [mdnslookup resolve](#mdnslookup-resolve)	 - Resolve a service instance to its host, port and attributes
* [mdnslookup types](#mdnslookup-types)	 - List the service types which are advertised on the link

## mdnslookup browse

Browse the instances of a service type

### Synopsis

Browse the instances of a service type, and report the instances as they are
added, updated and removed. The browse runs until it is interrupted, or until
the duration elapses when --duration is set. Without a service type, it browses
the service types which are advertised on the link.

A responder usually answers a browse with the host, the port and the addresses
of each instance, but one may answer only with the instance names. With
--resolve, such an instance is resolved before it is reported, as
avahi-browse --resolve does.

```
mdnslookup browse [service] [flags]
```

### Examples

```
  mdnslookup browse _matterc._udp
  mdnslookup browse --subtype _S3 _matterc._udp
  mdnslookup browse --resolve _http._tcp
  mdnslookup browse --duration 10s --format json _matter._tcp
  mdnslookup browse
```

### Options

```
      --domain string       domain to browse (default "local")
      --duration duration   browse duration (0 means until interrupted)
  -h, --help                help for browse
      --resolve             resolve an instance which is reported without its host, port or addresses
      --subtype string      service subtype to browse (_S3, _L840, ...)
      --unicast             request unicast responses (QU)
```

### Options inherited from parent commands

```
      --debug               enable debug output
      --family string       address family to use: all|ipv4|ipv6 (default "all")
      --format string       output format: table|json|csv (default "table")
  -i, --interface strings   network interfaces to use (all available interfaces by default)
  -t, --timeout duration    query timeout (default 5s)
      --verbose             enable verbose output
```

### SEE ALSO

* [mdnslookup](#mdnslookup)	 - Browse and resolve mDNS (DNS-SD) services on the local link

## mdnslookup decode

Decode a recorded mDNS message

### Synopsis

Decode a recorded mDNS message, and print its records. The file is a hex dump
log, such as the dumps which "monitor --hex" prints, or a raw message file when
--raw is set.

The command needs no network, so it can be used to analyze a message which was
captured elsewhere.

```
mdnslookup decode <file> [flags]
```

### Examples

```
  mdnslookup decode mdnstest/dumps/matter-answer-01.dump
  mdnslookup decode --format json mdnstest/dumps/matter-query-01.dump
  mdnslookup decode --raw message.bin
```

### Options

```
  -h, --help   help for decode
      --raw    read the file as a raw message instead of a hex dump log
```

### Options inherited from parent commands

```
      --debug               enable debug output
      --family string       address family to use: all|ipv4|ipv6 (default "all")
      --format string       output format: table|json|csv (default "table")
  -i, --interface strings   network interfaces to use (all available interfaces by default)
  -t, --timeout duration    query timeout (default 5s)
      --verbose             enable verbose output
```

### SEE ALSO

* [mdnslookup](#mdnslookup)	 - Browse and resolve mDNS (DNS-SD) services on the local link

## mdnslookup host

Resolve a host name to its addresses

### Synopsis

Resolve a host name to its addresses with Multicast DNS (RFC 6762).

The IPv6 scoped addressing zone is set for the link-local addresses, so that the
printed addresses can be used to connect to the host.

```
mdnslookup host <hostname> [flags]
```

### Examples

```
  mdnslookup host macmini.local
  mdnslookup host --family ipv6 84FCE6036F38.local
```

### Options

```
  -h, --help   help for host
```

### Options inherited from parent commands

```
      --debug               enable debug output
      --family string       address family to use: all|ipv4|ipv6 (default "all")
      --format string       output format: table|json|csv (default "table")
  -i, --interface strings   network interfaces to use (all available interfaces by default)
  -t, --timeout duration    query timeout (default 5s)
      --verbose             enable verbose output
```

### SEE ALSO

* [mdnslookup](#mdnslookup)	 - Browse and resolve mDNS (DNS-SD) services on the local link

## mdnslookup monitor

Watch the mDNS messages on the link

### Synopsis

Listen on the mDNS multicast address, and print every message which is sent
on the link. The monitor runs until it is interrupted.

```
mdnslookup monitor [flags]
```

### Examples

```
  mdnslookup monitor
  mdnslookup monitor --responses --format json
  mdnslookup monitor --hex
```

### Options

```
  -h, --help        help for monitor
      --hex         print the raw bytes of the messages
      --queries     print only the query messages
      --responses   print only the response messages
```

### Options inherited from parent commands

```
      --debug               enable debug output
      --family string       address family to use: all|ipv4|ipv6 (default "all")
      --format string       output format: table|json|csv (default "table")
  -i, --interface strings   network interfaces to use (all available interfaces by default)
  -t, --timeout duration    query timeout (default 5s)
      --verbose             enable verbose output
```

### SEE ALSO

* [mdnslookup](#mdnslookup)	 - Browse and resolve mDNS (DNS-SD) services on the local link

## mdnslookup query

Send a single question and print the answering records

### Synopsis

Send a single question to the multicast address, and print the records of
every answer until the timeout elapses.

Use it to look at the raw records of a responder. Use "browse" and "resolve" to
discover the services instead.

```
mdnslookup query [name] [flags]
```

### Examples

```
  mdnslookup query _matterc._udp.local
  mdnslookup query --type SRV DD200C20D25AE5F7._matterc._udp.local
  mdnslookup query --type ANY --format json macmini.local
```

### Options

```
  -h, --help          help for query
      --type string   question record type: PTR|SRV|TXT|A|AAAA|ANY (default "PTR")
      --unicast       request unicast responses (QU)
```

### Options inherited from parent commands

```
      --debug               enable debug output
      --family string       address family to use: all|ipv4|ipv6 (default "all")
      --format string       output format: table|json|csv (default "table")
  -i, --interface strings   network interfaces to use (all available interfaces by default)
  -t, --timeout duration    query timeout (default 5s)
      --verbose             enable verbose output
```

### SEE ALSO

* [mdnslookup](#mdnslookup)	 - Browse and resolve mDNS (DNS-SD) services on the local link

## mdnslookup resolve

Resolve a service instance to its host, port and attributes

### Synopsis

Resolve a service instance name to its host, port, addresses and TXT
attributes.

RFC 6763 (5) resolves an instance by querying its SRV and TXT records, and the
address records of the SRV target are queried when the responder does not send
them with the answer.

```
mdnslookup resolve <instance> [flags]
```

### Examples

```
  mdnslookup resolve DD200C20D25AE5F7._matterc._udp.local
  mdnslookup resolve --format json 95BDD2C0BDB33593._matterc._udp.local
```

### Options

```
  -h, --help   help for resolve
```

### Options inherited from parent commands

```
      --debug               enable debug output
      --family string       address family to use: all|ipv4|ipv6 (default "all")
      --format string       output format: table|json|csv (default "table")
  -i, --interface strings   network interfaces to use (all available interfaces by default)
  -t, --timeout duration    query timeout (default 5s)
      --verbose             enable verbose output
```

### SEE ALSO

* [mdnslookup](#mdnslookup)	 - Browse and resolve mDNS (DNS-SD) services on the local link

## mdnslookup types

List the service types which are advertised on the link

### Synopsis

List the service types which are advertised on the link.

RFC 6763 (7.2) enumerates the service types by browsing "_services._dns-sd._udp".

```
mdnslookup types [flags]
```

### Examples

```
  mdnslookup types
  mdnslookup types --timeout 10s --format json
```

### Options

```
      --domain string   domain to browse (default "local")
  -h, --help            help for types
```

### Options inherited from parent commands

```
      --debug               enable debug output
      --family string       address family to use: all|ipv4|ipv6 (default "all")
      --format string       output format: table|json|csv (default "table")
  -i, --interface strings   network interfaces to use (all available interfaces by default)
  -t, --timeout duration    query timeout (default 5s)
      --verbose             enable verbose output
```

### SEE ALSO

* [mdnslookup](#mdnslookup)	 - Browse and resolve mDNS (DNS-SD) services on the local link

