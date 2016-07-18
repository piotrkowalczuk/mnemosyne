package mnemosyne

import (
	"flag"

	_ "github.com/lib/pq"
)

var (
	testPostgresAddress string
)

func init() {
	flag.StringVar(&testPostgresAddress, "storage.postgres.address", "postgres://postgres:@localhost/test?sslmode=disable", "")
}
