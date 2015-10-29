package main

import (
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
)

var (
	storage Storage
)

func initStorage(fn func() (Storage, error), logger log.Logger) {
	s, err := fn()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	err = s.Setup()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	storage = s
}

func initPostgresStorage(tableName string, postgres *sql.DB) func() (Storage, error) {
	return func() (Storage, error) {
		return newPostgresStorage(postgres, tableName), nil
	}
}
