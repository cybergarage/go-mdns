# mdnsd

`mdnsd` is a Multicast DNS responder. It publishes one DNS-SD service on the link: it announces the service, answers the queries for it, and withdraws it with goodbye records when it is interrupted.

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

The host name resolves to the addresses of the interface each query arrives on. `mdnsd` does not probe for name conflicts yet, so choose a name no other node on the link uses.
