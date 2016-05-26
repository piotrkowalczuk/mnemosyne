package mnemosyne

import (
	"flag"

	_ "github.com/lib/pq"
)

var (
	testPostgresAddress string
)

func init() {
	flag.StringVar(&testPostgresAddress, "s.p.address", "postgres://postgres:@localhost/test?sslmode=disable", "")
}
