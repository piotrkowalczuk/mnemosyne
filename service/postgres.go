package service

import (
	"database/sql"
	"time"

	// ...
	"github.com/go-soa/mnemosyne/lib"
	_ "github.com/lib/pq"
)

// Postgres ...
var Postgres *sql.DB

// PostgresConfig ...
type PostgresConfig struct {
	Retry            bool   `xml:"retry"`
	TableName        string `xml:"table-name"`
	ConnectionString string `xml:"connection-string"`
}

// InitPostgres ...
func InitPostgres(config PostgresConfig, logger lib.StdLogger) {
	var err error

	// Because of recursion it needs to be checked to not spawn more than one.
	if Postgres == nil {
		Postgres, err = sql.Open("postgres", config.ConnectionString)
		if err != nil {
			logger.Fatal(err)
		}
	}

	// At this moment connection is not yet established.
	// Ping is required.
	if err := Postgres.Ping(); err != nil {
		if !config.Retry {
			logger.Fatal(err)
		}

		logger.Print(err)
		time.Sleep(2 * time.Second)

		InitPostgres(config, logger)
	} else {
		logger.Print("Connection do PostgreSQL established successfully.")
	}
}
