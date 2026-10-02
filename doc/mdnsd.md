# mdnsd

`mdnsd` is a Multicast DNS responder. It publishes one DNS-SD service on the link: it probes the instance name and the host name, announces the service, answers the queries for it, and withdraws it with goodbye records when it is interrupted.

Use [mdnslookup](mdnslookup.md) to browse and resolve the services.

## Synopsis

```
mdnsd -name <instance> -service <type> -port <port> [-host <host>] [-subtype <label>]... [-txt <string>]... [-v]
```

## Options

```
  -name string       service instance name
  -service string    service type, such as _http._tcp
  -port int          service port
  -host string       host name without the domain (default: this host)
  -subtype string    subtype label, such as _printer; can be repeated
  -txt string        TXT string, such as key=value; can be repeated
  -v                 print the received queries
```

Without `-service`, `mdnsd` publishes nothing and only prints the received queries with `-v`.

## Examples

```
$ mdnsd -name demo -service _http._tcp -port 8080 -txt path=/
```

A Matter commissionable node:

```
$ mdnsd -name 665F6E75B5D3A9C2 -service _matterc._udp -port 5540 -host B75AFB458ECD6D6F \
    -subtype _L3840 -subtype _S15 -subtype _CM -txt D=3840 -txt CM=1 -txt VP=65521+32769
```

The host name resolves to the addresses of the interface each query arrives on.

## Name conflicts

`mdnsd` does not rename the service as `dns-sd -R` and `avahi-publish-service` do. It exits with an error when another node on the link holds the instance name or the host name, when it starts or later:

```
$ mdnsd -name demo -service _http._tcp -port 8080
mdns: name conflict: demo._http._tcp.local
```

Run it again with another `-name` or `-host`.
