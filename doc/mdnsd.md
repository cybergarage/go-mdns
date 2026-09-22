# mdnsd

**mdnsd is under development.**

`mdnsd` is intended to become a Multicast DNS responder, which registers the services and answers the queries of the other nodes on the link. The responder side of go-mdns is not implemented yet, so the current command only listens for the mDNS messages and prints them when `-v` is set. It registers no service, and it answers no query.

Use [mdnslookup](mdnslookup.md) to browse and resolve the services. The responder, and this command with it, is planned for v1.0.0.

## Synopsis

```
mdnsd [-v]
```

## Options

```
  -v    enable verbose output
```
