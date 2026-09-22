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

This tool is a client. Advertising a service is not supported yet.

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

* [mdnslookup completion]()	 - Generate the autocompletion script for the specified shell
* [mdnslookup decode]()	 - Decode a recorded mDNS message
* [mdnslookup doc]()	 - Generate markdown documentation to stdout
* [mdnslookup host]()	 - Resolve a host name to its addresses
* [mdnslookup monitor]()	 - Watch the mDNS messages on the link
* [mdnslookup query]()	 - Send a single question and print the answering records
* [mdnslookup resolve]()	 - Resolve a service instance to its host, port and attributes
* [mdnslookup types]()	 - List the service types which are advertised on the link

## mdnslookup browse

Browse the instances of a service type

### Synopsis

Browse the instances of a service type, and report the instances as they are
added, updated and removed. The browse runs until it is interrupted, or until
the duration elapses when --duration is set.

```
mdnslookup browse [service] [flags]
```

### Examples

```
  mdnslookup browse _matterc._udp
  mdnslookup browse --subtype _S3 _matterc._udp
  mdnslookup browse --duration 10s --format json _matter._tcp
```

### Options

```
      --domain string       domain to browse (default "local")
      --duration duration   browse duration (0 means until interrupted)
  -h, --help                help for browse
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


## mdnslookup completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for mdnslookup for the specified shell.
See each sub-command's help for details on how to use the generated script.


### Options

```
  -h, --help   help for completion
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

* [mdnslookup completion bash]()	 - Generate the autocompletion script for bash
* [mdnslookup completion fish]()	 - Generate the autocompletion script for fish
* [mdnslookup completion powershell]()	 - Generate the autocompletion script for powershell
* [mdnslookup completion zsh]()	 - Generate the autocompletion script for zsh

## mdnslookup completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(mdnslookup completion bash)

To load completions for every new session, execute once:

#### Linux:

	mdnslookup completion bash > /etc/bash_completion.d/mdnslookup

#### macOS:

	mdnslookup completion bash > $(brew --prefix)/etc/bash_completion.d/mdnslookup

You will need to start a new shell for this setup to take effect.


```
mdnslookup completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
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


## mdnslookup completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	mdnslookup completion fish | source

To load completions for every new session, execute once:

	mdnslookup completion fish > ~/.config/fish/completions/mdnslookup.fish

You will need to start a new shell for this setup to take effect.


```
mdnslookup completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
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


## mdnslookup completion help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type completion help [path to command] for full details.

```
mdnslookup completion help [command] [flags]
```

### Options

```
  -h, --help   help for help
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


## mdnslookup completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	mdnslookup completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
mdnslookup completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
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


## mdnslookup completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(mdnslookup completion zsh)

To load completions for every new session, execute once:

#### Linux:

	mdnslookup completion zsh > "${fpath[1]}/_mdnslookup"

#### macOS:

	mdnslookup completion zsh > $(brew --prefix)/share/zsh/site-functions/_mdnslookup

You will need to start a new shell for this setup to take effect.


```
mdnslookup completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
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


## mdnslookup doc

Generate markdown documentation to stdout

```
mdnslookup doc [flags]
```

### Options

```
  -h, --help   help for doc
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


## mdnslookup help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type mdnslookup help [path to command] for full details.

```
mdnslookup help [command] [flags]
```

### Options

```
  -h, --help   help for help
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


