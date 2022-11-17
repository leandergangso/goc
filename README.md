# GOC

Simple CLI for tracking work hours into Google calendar.

## Setup (Linux)

- Clone this repo `git clone https://github.com/LeanderGangso/goc.git`
- Build the executable `go build cmd/goc.go`
- Make sure `goc` is executable, if not run: `chmod +x goc`
- Move the `goc` file into `/usr/local/bin` or any other folder that is in your `$PATH`
- Run `goc setup` to configure google app and calendar
- Run `goc` to see help and usage, and `goc help COMMAND` to see command info

## Setup (other)

You can use this both for Mac and Windwos, but tweaks may be needed.
Just follow the steps above and make changes where applicable.

## Google App

Note that, unless you put your google app in *production*, the token only lasts for 7 days.
If the refresh token expires the `token.json` will be deleted and you will be prompted to reauthenticate on he next command.

> You can set your application in production mode without a review, but will show a warning when users authenticate with your app.

## Usage Examples

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