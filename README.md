# goc

Simple CLI for tracking work hours into Google calendar.

## Setup (Linux)

- Download the latest `goc` file from the **release** page
- Make sure `goc` is executable, if not run: `chmod +x goc`
- Setup and download the Google [credentials.json](https://console.cloud.google.com/apis/credentials) file
  - Need to rename it to `credentials.json` after you have installed in from Google
- Move the `goc` and `credentials.json` files into `/usr/local/bin`
- Run `goc` to see help and usage, and `goc help COMMAND` to see command info
