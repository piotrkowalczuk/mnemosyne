package cli

import (
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/go-soa/mnemosyne/service"
)

var (
	initPostgresCommand = cli.Command{
		Name:   "initpostgres",
		Usage:  "set up postgres database",
		Action: initPostgresCommandAction,
	}
)

func initPostgresCommandAction(context *cli.Context) {
	service.InitConfig(context.GlobalString("environment"))
	service.InitLogger(service.Config.Logger)
	service.InitPostgres(service.Config.SessionStorage.Postgres)

	queryBytes, err := ioutil.ReadFile("data/sql/schema_postgres.sql")
	if err != nil {
		service.Logger.Fatal(err)
	}

	service.Logger.Info("Schema file opened successfully.")

	_, err = service.Postgres.Exec(string(queryBytes))
	if err != nil {
		service.Logger.Fatal(err)
	}

	service.Logger.Info("Postgres database has been prepared successfully.")
}
