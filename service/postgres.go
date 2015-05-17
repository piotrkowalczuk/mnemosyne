package service

import (
	"database/sql"
	"time"

	// ...
	_ "github.com/lib/pq"
)

// Postgres ...
var Postgres *sql.DB

// PostgresConfig ...
type PostgresConfig struct {
	TableName        string `xml:"table-name"`
	ConnectionString string `xml:"connection-string"`
}

// InitPostgres ...
func InitPostgres(config PostgresConfig) {
	var err error

	// Because of recursion it needs to be checked to not spawn more than one.
	if Postgres == nil {
		Postgres, err = sql.Open("postgres", config.ConnectionString)
		if err != nil {
			Logger.Fatal(err)
		}
	}

	// At this moment connection is not yet established.
	// Ping is required.
	if err := Postgres.Ping(); err != nil {
		Logger.Error(err)
		time.Sleep(2 * time.Second)

		InitPostgres(config)
	} else {
		Logger.Info("Connection do PostgreSQL established successfully.")
	}
}
