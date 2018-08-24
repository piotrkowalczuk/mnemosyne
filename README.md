# Mnemosyne [![CircleCI](https://circleci.com/gh/piotrkowalczuk/mnemosyne.svg?style=svg)](https://circleci.com/gh/piotrkowalczuk/mnemosyne)

[![GoDoc](https://godoc.org/github.com/piotrkowalczuk/mnemosyne?status.svg)](http://godoc.org/github.com/piotrkowalczuk/mnemosyne)
[![Docker Pulls](https://img.shields.io/docker/pulls/piotrkowalczuk/mnemosyne.svg?maxAge=604800)](https://hub.docker.com/r/piotrkowalczuk/mnemosyne/)
[![codecov.io](https://codecov.io/github/piotrkowalczuk/mnemosyne/coverage.svg?branch=master)](https://codecov.io/github/piotrkowalczuk/mnemosyne?branch=master)
[![Code Climate](https://codeclimate.com/github/piotrkowalczuk/mnemosyne/badges/gpa.svg)](https://codeclimate.com/github/piotrkowalczuk/mnemosyne)
[![pypi](https://img.shields.io/pypi/v/mnemosyne-client.svg)](https://pypi.python.org/pypi/mnemosyne-client)

## Introduction

[Mnemosyne](http://github.com/piotrkowalczuk/mnemosyne) is an open-source self-hosted session management service.
It's written in Go, making it easy to build and deploy as a static binary.

It provides gRPC [interface](https://godoc.org/github.com/piotrkowalczuk/mnemosyne/mnemosynerpc#SessionManager). 
Messages are encoded using protobuf.

### Quick Start

To install and run service:

```bash
$ go get -d github.com/piotrkowalczuk/mnemosyne/...
$ cd $GOPATH/src/github.com/piotrkowalczuk/mnemosyne
$ make
$ mnemosyned -log.environment=development -postgres.address='postgres://localhost/example?sslmode=disable'
```

### Storage Engine
Goal is to support multiple storage's, like [PostgreSQL](http://www.postgresql.org/), [Redis](http://redis.io) or [MongoDB](https://www.mongodb.org). 
Nevertheless currently supported is only [PostgreSQL](http://www.postgresql.org/).

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
It is worth to mention that `latest` tag is released after each successful master branch build. Please use only images tagged using specific version anywhere else than a local development environment.

### From source

To install from source both go tools and [dep](https://github.com/golang/dep) is required. 

```
$ go get -d github.com/piotrkowalczuk/mnemosyne/...
$ cd $GOPATH/src/github.com/piotrkowalczuk/mnemosyne
$ make
```

### Configuration
`mnemosyned` accepts command line arguments to control its behavior. 
Possible options are listed below.

| Name | Flag | Default | Type |
| --- | --- | --- | --- |
| host | `-host` | 127.0.0.1 | string |
| port | `-port` | 8080 | int |
| grpc debug mode| `-grpc.debug` | false | boolean |
| cluster listen address | `-cluster.listen` | | string |
| cluster seeds | `-cluster.seeds` | | string |
| time to live | `-ttl` | 24m | duration |
| time to clear | `-ttc` | 1m | duration |
| logger environment | `-log.environment` | production | enum(development, production, stackdriver) |
| logger level | `-log.level` | info | enum(debug, info, warn, error, dpanic, panic, fatal) |
| storage | `-storage` | postgres | enum(postgres) |
| postgres address | `-postgres.address` | postgres://postgres:postgres@postgres/postgres?sslmode=disable | string |
| postgres table | `-postgres.table` | session | string |
| postgres schema | `-postgres.schema` | mnemosyne | string |
| tls | `-tls` | false | boolean |
| tls certificate file | `-tls.crt` | | string |
| tls key file |`-tls.key` | | string |

### Running

As we know, mnemosyne can be configured in many ways. For the beginning we can start simple:

```bash
$ mnemosyned postgres.address="postgres://localhost/test?sslmode=disable"
```
Mnemosyne will automatically create all required tables/indexes for specified database.

### Monitoring
`mnemosyned` works well with [Prometheus](http://prometheus.io). 
It exposes multiple metrics through `/metrics` endpoint, it includes:

* `mnemosyned_cache_hits_total`
* `mnemosyned_cache_misses_total`
* `mnemosyned_cache_refresh_total`
* `mnemosyned_storage_postgres_errors_total`
* `mnemosyned_storage_postgres_queries_total`
* `mnemosyned_storage_postgres_query_duration_seconds`
* `mnemosyned_storage_postgres_connections`

Additionally to that `mnemosyned` is using internally [promgrpc](https://github.com/piotrkowalczuk/promgrpc) package to monitor entire incoming and outgoing RPC traffic.

### Examples

#### Go

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

#### Python

Library is available through [pypi](https://pypi.python.org/pypi/mnemosyne-client) and can be installed by typing `pip install mnemosyne-client`.

```python
from  mnemosynerpc import session_pb2, session_pb2_grpc
import grpc


channel = grpc.insecure_channel('localhost:8080')
stub = session_pb2_grpc.SessionManagerStub(channel)

for i in range(0, 10):
	res = stub.Start(session_pb2.StartRequest(session=session_pb2.Session(subject_id=str(i))))

	res = stub.Get(session_pb2.GetRequest(access_token=res.session.access_token))
	print "%s - %s" % (res.session.access_token, res.session.expire_at.ToJsonString())
```

## Contribution

TODO: describe