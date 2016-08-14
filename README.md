# Mnemosyne [![Build Status](https://travis-ci.org/piotrkowalczuk/mnemosyne.svg)](https://travis-ci.org/piotrkowalczuk/mnemosyne)

[![GoDoc](https://godoc.org/github.com/piotrkowalczuk/mnemosyne?status.svg)](http://godoc.org/github.com/piotrkowalczuk/mnemosyne)
[![Docker Pulls](https://img.shields.io/docker/pulls/piotrkowalczuk/mnemosyne.svg?maxAge=604800)](https://hub.docker.com/r/piotrkowalczuk/mnemosyne/)
[![codecov.io](https://codecov.io/github/piotrkowalczuk/mnemosyne/coverage.svg?branch=master)](https://codecov.io/github/piotrkowalczuk/mnemosyne?branch=master)
[![Code Climate](https://codeclimate.com/github/piotrkowalczuk/mnemosyne/badges/gpa.svg)](https://codeclimate.com/github/piotrkowalczuk/mnemosyne)
[![Gitter](https://badges.gitter.im/piotrkowalczuk/mnemosyne.svg)](https://gitter.im/piotrkowalczuk/mnemosyne?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Introduction

[Mnemosyne](http://github.com/piotrkowalczuk/mnemosyne) is an open-source self-hosted session management service.
It's written in Go, making it easy to build and deploy as a static binary.

It provides two ways for communication:

* [mnemosyne.Mnemosyne](https://godoc.org/github.com/piotrkowalczuk/mnemosyne#Mnemosyne) - Simplified way that hides complexity of [gRPC](http://www.grpc.io) library.
* [mnemosynerpc.SessionManager](https://godoc.org/github.com/piotrkowalczuk/mnemosyne/mnemosynerpc#SessionManager) - Full feature client.

### Quick Start

To install and run service:

```bash
$ go get -d github.com/piotrkowalczuk/mnemosyne/...
$ cd $GOPATH/src/github.com/piotrkowalczuk/mnemosyne
$ glide install
$ go install ./cmd/mnemosyned
$ mnemosyned -log.format=humane -postgres.address='postgres://localhost/example?sslmode=disable'
```

Simpliest implementation could looks like that:

```go
package main

import (
	"fmt"

	"golang.org/x/net/context"
	"github.com/piotrkowalczuk/mnemosyne"
)

func main() {
	mnemo, err := mnemosyne.New(mnemosyne.MnemosyneOpts{
		Addresses: []string{"127.0.0.1:8080"},
		Block: true,
	})
	if err != nil {
		// ...
	}
	defer mnemo.Close()

	ses, err := mnemo.Start(context.Background(), "subject-id", "subject-client", map[string]string{
		"username": "johnsnow@gmail.com",
		"first_name": "John",
		"last_name": "Snow",
	})
	if err != nil {
		// ...
	}

	fmt.Println(ses.AccessToken)
}
```
### Storage Engine
Goal is to support multiple storage's, like [PostgreSQL](http://www.postgresql.org/), [Redis](http://redis.io) or [MongoDB](https://www.mongodb.org). Nevertheless currently supported is only [PostgreSQL](http://www.postgresql.org/).

### Remote Procedure Call API
For communication, Mnemosyne is exposing RPC API that uses [protocol buffers](https://developers.google.com/protocol-buffers/), Google’s mature open source mechanism for serializing structured data.

* Create
* Get
* List
* Exists
* Abandon
* SetData
* Delete

## Installation

Mnemosyne can be installed in one way, from source.
Or can be used as a container using docker [image](https://hub.docker.com/r/piotrkowalczuk/mnemosyne/).

### From source

To install from source both go tools and [glide](github.com/Masterminds/glide) is required. 

```
$ go get -d github.com/piotrkowalczuk/mnemosyne/...
$ cd $GOPATH/src/github.com/piotrkowalczuk/mnemosyne
$ glide install
$ go install ./cmd/mnemosyned
```

### Configuration
**mnemosyned** accepts command line arguments to control its behavior. Possible options are is listed below.

| Name | Flag | Default | Type |
| --- | --- | --- | --- |
| host | `-host` | 127.0.0.1 | string |
| port | `-port` | 8080 | int |
| time to live | `-ttl` | 24m | duration |
| time to clear | `-ttc` | 1m | duration |
| logger format | `-log.format` | json | enum(json, humane, logfmt) |
| logger adapter | `-log.adapter` | stdout | enum(stdout) |
| monitoring | `-monitoring ` | false | boolean |
| storage | `-storage` | postgres | enum(postgres) |
| postgres address | `-postgres.address` | postgres://postgres:postgres@postgres/postgres?sslmode=disable | string |
| tls | `-tls` | false | boolean |
| tls certificate file | `-tls.cert` | | string |
| tls key file |`-tls.key` | | string |

### Running

As we know, mnemosyne can be configured in many ways. For the beginning we can start simple:

```bash
$ mnemosyned postgres.address="postgres://localhost/test?sslmode=disable"
```

Mnemosyne will automatically create all required tables/indexes for specified database.

## Contribution

### TODO

- [ ] Client library
    - [x] Go
    - [ ] Python
- [ ] Engines
	- [x] PostgreSQL
		- [x] Get
		- [x] List
		- [x] Exists
		- [x] Create
		- [x] Abandon
		- [x] SetData
		- [x] Delete
		- [x] Setup
		- [x] TearDown
	- [ ] RAM
	- [ ] Redis