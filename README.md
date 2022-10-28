# GOC

Simple CLI for tracking work hours into Google calendar.

## Setup (Linux)

- Download the latest `goc` file from the **release** page or clone this repo and build it yourself
- Make sure `goc` is executable, if not run: `chmod +x goc`
- Move the `goc` file into `/usr/local/bin` or any other folder that is in your `$PATH`
- Run `goc setup` to configure the calendar.
- Run `goc` to see help and usage, and `goc help COMMAND` to see command info

## Setup (Other)

If you want to use this on Mac and Windows, you need to make some changes.
This includes, but may not be limited to the following:

- You need to build the executable on your own, the one in the **release** page will not work
  - Need `go` installed on your system
  - Run `go build cmd/goc.go` to build the executable
- The `$HOME` environment variable is used, if your system don't support this, you can hardcode a path instead of $HOME

## Usage examples

See help:
```bash
// show help
goc
// show command help
goc help start
// show help for start command
goc help start
```

Basic usage:
```bash
// start task at the current time
goc start 'task name'
// see status of current task
goc status
// end current task at the current time
goc end
```

Custom times (format: HH:MM):
```bash
// start task at a different time
goc start 'task name' -t 8:00
// start new task that will end the previous task
goc start 'new task' -t 10:00
// end current task at a different time
goc end -t 16:00
```

Alias usage:
```bash
// list aliases
goc list
// new alias
goc alias 'alias name' 'task name'
// use alias
goc start 'alias name'
// remove alias
goc remove 'alias name'
```