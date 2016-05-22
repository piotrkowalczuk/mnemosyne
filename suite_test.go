package mnemosyne

import "flag"

var (
	testPostgresAddress string
)

func init() {
	flag.StringVar(&testPostgresAddress, "s.p.address", "postgres://postgres:@localhost/test?sslmode=disable", "")
}
