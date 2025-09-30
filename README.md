# rep-cmd

🌀 **rep-cmd** is a simple CLI tool that repeatedly executes any command at a fixed interval.  
Unlike `cron`, it requires no configuration files—just start the process and stop it whenever you want.  
To stop, simply press `Ctrl+C` or send a `kill` signal.

## Installation

```bash
go install github.com/hxrxchang/rep-cmd@latest
```

After installation, the binary will be placed in `$GOPATH/bin` or `$HOME/go/bin`.

## Usage

```bash
rep-cmd -c 'your command' -i interval
```

- `-c` : The command to run (**required**)
- `-i` : Interval between executions. Supports `10s`, `5m`, `1h` (default: `5m`)

### Example: Open Slack web app every 5 minutes

```bash
rep-cmd -c 'open -a "Google Chrome" https://slack.com/app' -i 5m
```
