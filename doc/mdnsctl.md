## mdnsctl



### Options

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
  -h, --help            help for mdnsctl
      --verbose         enable verbose output
```

* [mdnsctl doc]()	 - Generate markdown documentation to stdout
* [mdnsctl dump]()	 - Dump mDNS messages.
* [mdnsctl scan]()	 - Scan for mDNS devices.

## mdnsctl completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for mdnsctl for the specified shell.
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

* [mdnsctl completion bash]()	 - Generate the autocompletion script for bash
* [mdnsctl completion fish]()	 - Generate the autocompletion script for fish
* [mdnsctl completion powershell]()	 - Generate the autocompletion script for powershell
* [mdnsctl completion zsh]()	 - Generate the autocompletion script for zsh

## mdnsctl completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(mdnsctl completion bash)

To load completions for every new session, execute once:

#### Linux:

	mdnsctl completion bash > /etc/bash_completion.d/mdnsctl

#### macOS:

	mdnsctl completion bash > $(brew --prefix)/etc/bash_completion.d/mdnsctl

You will need to start a new shell for this setup to take effect.


```
mdnsctl completion bash
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


## mdnsctl completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	mdnsctl completion fish | source

To load completions for every new session, execute once:

	mdnsctl completion fish > ~/.config/fish/completions/mdnsctl.fish

You will need to start a new shell for this setup to take effect.


```
mdnsctl completion fish [flags]
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


## mdnsctl completion help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type completion help [path to command] for full details.

```
mdnsctl completion help [command] [flags]
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


## mdnsctl completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	mdnsctl completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
mdnsctl completion powershell [flags]
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


## mdnsctl completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(mdnsctl completion zsh)

To load completions for every new session, execute once:

#### Linux:

	mdnsctl completion zsh > "${fpath[1]}/_mdnsctl"

#### macOS:

	mdnsctl completion zsh > $(brew --prefix)/share/zsh/site-functions/_mdnsctl

You will need to start a new shell for this setup to take effect.


```
mdnsctl completion zsh [flags]
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


## mdnsctl doc

Generate markdown documentation to stdout

```
mdnsctl doc [flags]
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


## mdnsctl dump

Dump mDNS messages.

### Synopsis

Dump mDNS messages.

```
mdnsctl dump [flags]
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


## mdnsctl help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type mdnsctl help [path to command] for full details.

```
mdnsctl help [command] [flags]
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


## mdnsctl scan

Scan for mDNS devices.

### Synopsis

Scan for mDNS devices.

```
mdnsctl scan [flags]
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


