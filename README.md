# osmshortlink-go: Go module and command-line tool for creating, encoding & decoding of OpenStreetMap shortlinks

[![Go Reference](https://pkg.go.dev/badge/github.com/stefanb/osmshortlink-go.svg)](https://pkg.go.dev/github.com/stefanb/osmshortlink-go)
[![Test](https://github.com/stefanb/osmshortlink-go/actions/workflows/test.yml/badge.svg)](https://github.com/stefanb/osmshortlink-go/actions/workflows/test.yml)
[![golangci-lint](https://github.com/stefanb/osmshortlink-go/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/stefanb/osmshortlink-go/actions/workflows/golangci-lint.yml)
[![CodeQL](https://github.com/stefanb/osmshortlink-go/actions/workflows/codeql.yml/badge.svg)](https://github.com/stefanb/osmshortlink-go/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/stefanb/osmshortlink-go)](https://goreportcard.com/report/github.com/stefanb/osmshortlink-go)
[![codebeat badge](https://codebeat.co/badges/0dcfa9c5-a59b-46ed-b0a6-30e1bbda9a7e)](https://codebeat.co/projects/github-com-stefanb-osmshortlink-go-main)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/stefanb/osmshortlink-go/badge)](https://securityscorecards.dev/viewer/?uri=github.com/stefanb/osmshortlink-go)
[![GitHub Release](https://img.shields.io/github/release/stefanb/osmshortlink-go.svg?style=flat)](https://github.com/stefanb/osmshortlink-go/releases/latest)

Specification: https://wiki.openstreetmap.org/wiki/Shortlink

## Usage

### Creating a link in Go

```go
package main

import "github.com/stefanb/osmshortlink-go"

func main() {
	shortlink, err := osmshortlink.Create(46.05141, 14.50604, 17)
	if err != nil {
		panic(err)
	}
	print(shortlink)
}
```

Prints: [`https://osm.org/go/0Ik3VNr_A-?m`](https://osm.org/go/0Ik3VNr_A-?m)

[Try it in Go playground](https://go.dev/play/p/mObcbRyGU9E)

### Command-line tool

#### Installation

##### Download the pre-built binaries

Download the pre-built binaries for your platform from [latest release](https://github.com/stefanb/osmshortlink-go/releases/latest).

##### Install using Go

````bash
$ go install github.com/stefanb/osmshortlink-go/cmd/osmshortlink@latest
$ osmshortlink 46.05141 14.50604 17
https://osm.org/go/0Ik3VNr_A-?m
````

#### Run using Go without installing

````bash
$ go run github.com/stefanb/osmshortlink-go/cmd/osmshortlink@latest 46.05141 14.50604 17
https://osm.org/go/0Ik3VNr_A-?m
````

#### Pre-built binaries

You can download pre-built binaries for various platforms from [latest release](https://github.com/stefanb/osmshortlink-go/releases/latest).

```bash
Usage: osmshortlink [latitude] [longitude] [zoom]
```

For example:

```bash
$ osmshortlink 46.05141 14.50604 17
https://osm.org/go/0Ik3VNr_A-?m
```

## Development

Run: `go run ./cmd/osmshortlink/main.go 46.05141 14.50604 17`

Testing: `go test -v -race -cover ./...`

Fuzz testing: `go test -fuzz=. -fuzztime=1m`
