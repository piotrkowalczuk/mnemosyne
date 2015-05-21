package cli

import (
	"os"

	"github.com/codegangsta/cli"
)

func init() {
	app := cli.NewApp()
	app.Name = "mnemosyne"
	app.Usage = "..."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "environment, e",
			Value:  "development",
			Usage:  "environment in wich application is running",
			EnvVar: "MNEMOSYNE_ENV",
		},
	}
	app.Commands = []cli.Command{
		runCommand,
		initCommand,
	}

	app.Run(os.Args)
}
