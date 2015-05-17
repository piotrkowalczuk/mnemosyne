Mnemosyne [![Build Status](https://travis-ci.org/go-soa/mnemosyne)](https://travis-ci.org/go-soa/mnemosyne)
=============
Installation
------------
1. Set you GOPATH properly (http://golang.org/doc/code.html#GOPATH)
2. `go get github.com/go-soa/mnemosyne`
3. `go get` if some dependencies are missing
4. Create `conf/{env}.xml` based on `conf/{env}.xml.dist`
5. Set `$MNEMOSYNE_ENV` global variable to `test`, `development` or `production`

Commands
--------

#### Build
```bash
go build
```

#### Service
```bash
./mnemosyne initpostgres - execute data/sql/schema_postgres.sql against configured database.
./mnemosyne run - starts server.
./mnemosyne help [command] - display help message about available commands
```

Dependencies
------------
- PostgreSQL

TODO
----
- [ ] Commands
	- [x] Initialize postgres database
	- [x] Start server
- [ ] Client library
    - [ ] Go
    - [ ] Python
- [ ] Engines
	- [x] PostgreSQL
		- [x] Get(SessionID) (*Session, error)
		- [x] Exists(SessionID) (bool, error)
		- [x] New(SessionData) (*Session, error)
		- [x] Abandon(SessionID) error
		- [x] SetData(SessionDataEntry) (*Session, error)
	- [ ] RAM
	- [ ] Redis
	- [ ] MySQL
	- [ ] MongoDB