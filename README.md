# GOC

Simple CLI for tracking work hours into Google calendar.

## Setup (Linux)

- Clone this repo and build the executable with `go build cmd/goc.go`
- Make sure `goc` is executable, if not run: `chmod +x goc`
- Move the `goc` file into `/usr/local/bin` or any other folder that is in your `$PATH`
- Run `goc setup` to configure google app and calendar
- Run `goc` to see help and usage, and `goc help COMMAND` to see command info

## Usage examples (may differ)

See help:
```bash
# show help
goc
# show command help
goc s -h
```

Basic usage:
```bash
# start task at the current time
goc s 'task name'
# see status of current task
goc st
# end current task at the current time
goc e
```

Custom times (format: HHMM):
```bash
# start task at a different time
goc s -t 0800 'task name'
# start new task that will end the previous task
goc s -t 1000 'new task'
# end current task at a different time
goc e -t 1600
```

Alias usage:
```bash
# list aliases
goc l
# new alias
goc a 'alias name' 'task name'
# use alias
goc s 'alias name' 'optional description'
# remove alias
goc r 'alias name'
```