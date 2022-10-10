# Google Calendar CLI

Simple CLI for tracking work hours into Google calendar.

## Setup (Linux)

- Download the latest `goc` file from the **release** page
- Make sure `goc` is executable, if not run: `chmod +x goc`
- Move the `goc` file into `/usr/local/bin`
- Setup and download the Google [credentials.json](https://console.cloud.google.com/apis/credentials) file
  - Click `create credentials` and select `OAuth client ID`
  - Set `Application type` as `Desktop app` and follow the steps
  - Choose `download JSON` after creating credential
  - Rename the file you downloaded to `credentials.json`
  - Move this file into `$HOME/.goc_cli`
- Reset(close/reopen) terminal window for changes to take effect
- Run `goc` to see help and usage, and `goc help COMMAND` to see command info

## Setup (Other)

If you want to use this on Max and Windows, you need to make some changes.
This includes, but may not be limited to the following:

- You need to build the executable on your own, the one in the **release** page will not work
  - Need `go` installed on your system
  - Run `go build cmd/goc.go` to build the executable
- The `$HOME` environment variable is used, if your system don't support this, you can hardcode a path instead of $HOME
