# Mnemosyne [![Build Status](https://travis-ci.org/go-soa/mnemosyne.svg)](https://travis-ci.org/go-soa/mnemosyne)

## Installation

1. Set you GOPATH properly (http://golang.org/doc/code.html#GOPATH)
2. `go get github.com/go-soa/mnemosyne`
3. `go get` if some dependencies are missing
4. Set environmental variables


## Environment Variables 

* MNEMOSYNE_HOST
* MNEMOSYNE_PORT
* MNEMOSYNE_SUBSYSTEM
* MNEMOSYNE_LOGGER_FORMAT
* MNEMOSYNE_LOGGER_ADAPTER
* MNEMOSYNE_LOGGER_LEVEL
* MNEMOSYNE_STORAGE_ENGINE
* MNEMOSYNE_STORAGE_POSTGRES_CONNECTION_STRING
* MNEMOSYNE_STORAGE_POSTGRES_TABLE_NAME
* MNEMOSYNE_STORAGE_POSTGRES_RETRY

## Commands

* `make build` - builds daemon application
* `make test` - starts all possible tests
*

## Dependencies

- PostgreSQL

##TODO

- [ ] Client library
    - [ ] Go
    - [ ] Python
- [ ] Engines
	- [x] PostgreSQL
		- [x] Get
		- [ ] List
		- [x] Exists
		- [x] Create
		- [x] Abandon
		- [x] SetData
		- [x] Delete
	- [ ] RAM
	- [ ] Redis
