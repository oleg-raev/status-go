
# Status bindings for go-ethereum [![Master Build Status](https://img.shields.io/travis/status-im/status-go/master.svg?label=build/master)](https://github.com/status-im/status-go/tree/master) [![Develop Build Status](https://img.shields.io/travis/status-im/status-go/develop.svg?label=build/develop)](https://github.com/status-im/status-go/tree/develop)

[![TravisCI Builds](https://img.shields.io/badge/TravisCI-URL-yellowgreen.svg?link=https://travis-ci.org/status-im/status-go)](https://travis-ci.org/status-im/status-go)
[![GoDoc](https://godoc.org/github.com/status-im/status-go?status.svg)](https://godoc.org/github.com/status-im/status-go)




# Docs

- [How To Build](https://github.com/status-im/status-go/wiki/Build-Process-Explained)
- [Notes on Bindings](https://github.com/status-im/status-go/wiki/Notes-on-Bindings)
- [Developers' Guide](https://github.com/status-im/status-react/wiki/Developers'-Guide)

# Intro
status-go is a underlying part of status mobile app and handles blockchain and Status accounts related code.

It's written in Go and requires Go 1.8 or above.

It uses Makefile to do most common actions. See `make help` output for available commands.

Integration tests and server require `statusd-data` dir to store synched blockchain and other stuff, so we use `build/env.sh` script before running most of the commands in Makefile (like `go test` or running the server)

# Build
There are two main modes status-go can be built:
 - standalone server
 - library to link for Android or iOS

Use following Makefile commands:

`make statusgo` (builds binary into `build/bin/statusd`)
`make statusgo-android`) (builds .aar file `build/android-16/aar`)
`make statusgo-ios` and `make statusgo-ios-simulator` (builds iOS related artifacts in `build/os-9.3/framework`)

# Testing
To test statusgo, use:
`make ci`

If you want to launch specific test - `RPCSendTransactions` for example - use following command:
`./build/env.sh go test -v ./geth/api/ -testify.m RPCSendTransaction$`

# Licence

[Mozilla Public License 2.0](https://github.com/status-im/status-go/blob/develop/LICENSE.md)