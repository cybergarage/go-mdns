## mdnslookup



### Options

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
  -h, --help            help for mdnslookup
      --verbose         enable verbose output
```

* [mdnslookup doc]()	 - Generate markdown documentation to stdout
* [mdnslookup dump]()	 - Dump mDNS messages.
* [mdnslookup query]()	 - Query for mDNS devices.
* [mdnslookup scan]()	 - Scan for mDNS devices.

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
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
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
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
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
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
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
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
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
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
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
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
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
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## mdnslookup dump

Dump mDNS messages.

### Synopsis

Dump mDNS messages.

```
mdnslookup dump [flags]
```

### Options

```
  -h, --help   help for dump
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
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
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## mdnslookup query

Query for mDNS devices.

### Synopsis

Query for mDNS devices.

```
mdnslookup query [service] [flags]
```

### Examples

```
query _matterc._udp.local
```

### Options

```
  -h, --help   help for query
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## mdnslookup scan

Scan for mDNS devices.

### Synopsis

Scan for mDNS devices packets.

```
mdnslookup scan [flags]
```

### Options

```
  -h, --help   help for scan
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


