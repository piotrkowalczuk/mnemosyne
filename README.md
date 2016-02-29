# Mnemosyne [![GoDoc](https://godoc.org/github.com/piotrkowalczuk/mnemosyne?status.svg)](http://godoc.org/github.com/piotrkowalczuk/mnemosyne) [![Build Status](https://travis-ci.org/piotrkowalczuk/mnemosyne.svg)](https://travis-ci.org/piotrkowalczuk/mnemosyne)&nbsp;[![Code Climate](https://codeclimate.com/github/piotrkowalczuk/mnemosyne/badges/gpa.svg)](https://codeclimate.com/github/piotrkowalczuk/mnemosyne)

## Documentation
Documentation is available on [mnemosyne.readme.io](http://mnemosyne.readme.io).

## Contribution

## TODO

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

## Building

Increment version in `mnemosynd/config.go`. Execute `make package`.

Changes to flags or flag value defaults should be into 
`scripts/mnemosyne.service` and `mnemosyne.env`.