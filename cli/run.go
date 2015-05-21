package cli

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/codegangsta/cli"
	mnemosynerpc "github.com/go-soa/mnemosyne/controller/rpc"
	"github.com/go-soa/mnemosyne/service"
)

var (
	runCommand = cli.Command{
		Name:   "run",
		Usage:  "starts server",
		Action: runCommandAction,
	}
)

func runCommandAction(context *cli.Context) {
	service.InitConfig("conf", context.GlobalString("environment"))
	service.InitLogger(service.Config.Logger)
	service.InitPostgres(service.Config.SessionStorage.Postgres, service.Logger)
	service.InitSessionStorage(service.Config.SessionStorage, service.Logger)

	server := rpc.NewServer()
	server.Register(&mnemosynerpc.Store{
		Storage: service.SessionStorage,
		Logger:  service.Logger,
	})

	listenOn := service.Config.Server.Host + ":" + service.Config.Server.Port
	listen, err := net.Listen("tcp", listenOn)
	if err != nil {
		service.Logger.Fatal(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			service.Logger.Fatal(err)
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
