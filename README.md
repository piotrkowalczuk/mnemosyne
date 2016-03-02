# Mnemosyne [![GoDoc](https://godoc.org/github.com/piotrkowalczuk/mnemosyne?status.svg)](http://godoc.org/github.com/piotrkowalczuk/mnemosyne) [![Build Status](https://travis-ci.org/piotrkowalczuk/mnemosyne.svg)](https://travis-ci.org/piotrkowalczuk/mnemosyne)&nbsp;[![Code Climate](https://codeclimate.com/github/piotrkowalczuk/mnemosyne/badges/gpa.svg)](https://codeclimate.com/github/piotrkowalczuk/mnemosyne)

## Introduction

[Mnemosyne](http://github.com/piotrkowalczuk/mnemosyne) is an open-source self-hosted session management service. It's written in Go, making it easy to build and deploy as a static binary.

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

Mnemosyne can be installed in two ways, from source and using `deb` package that can be found in dist directory.

### From source

You can directly use the go tool to download and install the **mnemosyned** binary into your [GOPATH](https://github.com/golang/go/wiki/GOPATH). Go in version 1.6 is required.

```
$ go install github.com/piotrkowalczuk/mnemosyne/mnemosyned
```

### Configuration
**mnemosyned** accepts command line arguments to control its behavior. Possible options are is listed below.

| Name | Flag | Default | Type |
| --- | --- | --- | --- |
| host | `-host` | 127.0.0.1 | string |
| port | `-port` | 8080 |int |
| logger format | `-l.format` | json | enum(json, humane, logfmt) |
| logger adapter | `-l.adapter` | stdout | enum(stdout) |
|namespace|`-namespace`|string||
|subsystem|`-subsystem`| mnemosyne|string|
|monitoring engine|`-m.engine`|prometheus|enum(prometheus)|
|storage engine|`-s.engine`|postgres|enum(postgres)|
|storage postgres connection string|`-sp.connectionstring`|postgres://localhost:5432?sslmode=disable|string|
|storage postgres table name|`-sp.tablename`|mnemosyne_session|string|

### Running

As we know, mnemosyne can be configured in many ways. For the beginning we can start simple:

```bash
$ mnemosyned -namespace=acme -sp.connectionstring="postgres://localhost/test?sslmode=disable"
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

### Building

Increment version in `mnemosynd/config.go`. Execute `make package`.

Changes to flags or flag value defaults should be into 
`scripts/mnemosyne.service` and `mnemosyne.env`.