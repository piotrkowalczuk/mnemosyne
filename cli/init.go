package cli

import (
	"github.com/codegangsta/cli"
	"github.com/go-soa/mnemosyne/lib"
	"github.com/go-soa/mnemosyne/service"
)

var (
	initCommand = cli.Command{
		Name:   "init",
		Usage:  "...",
		Action: initCommandAction,
	}
)

func initCommandAction(context *cli.Context) {
	service.InitConfig("conf", context.GlobalString("environment"))

	ssc := service.Config.SessionStorage

	service.InitLogger(service.Config.Logger)

	switch ssc.Engine {
	case lib.PostgresEngine:
		service.InitPostgres(ssc.Postgres, service.Logger)
		ps := lib.NewPostgresStorage(service.Postgres, ssc.Postgres.TableName)

		err := ps.Init()
		if err != nil {
			service.Logger.Fatal(err)
		}

		service.Logger.Info("Postgres database has been prepared successfully.")
	}

	service.Logger.Info("Environment has been prepared successfully.")
}
