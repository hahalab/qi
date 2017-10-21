package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/todaychiji/ha/conf"
)

func main() {
	app := cli.NewApp()
	app.Name = "up"
	app.Version = version
	app.Usage = "build and deploy any time any where"
	app.Action = initServer
	app.Flags = conf.Flags
	//app.Before = conf.MustParseConfig
	app.Commands = []cli.Command{
		{
			Name:     "build",
			Usage:    "start server",
			Flags:    conf.Flags,
			Action:   initServer,
			Category: "BUILD",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var version string

func initServer(c *cli.Context) error {
	fmt.Println(c.String(conf.FlagConfigFile))
	return nil
}
