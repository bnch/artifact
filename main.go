package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	a := cli.NewApp()

	a.Name = "artifact"
	a.Usage = "make build artifacts beautifully"
	a.Action = httpServer
	a.Commands = []cli.Command{
		{
			Name:   "add-repo",
			Usage:  "add a new github repo to artifact",
			Action: addRepo,
		},
		{
			Name:   "mkconf",
			Usage:  "create the config file (conf.json)",
			Action: mkconf,
		},
		{
			Name:   "mkdb",
			Usage:  "automigrate the database",
			Action: mkdb,
		},
	}

	a.Run(os.Args)
}
