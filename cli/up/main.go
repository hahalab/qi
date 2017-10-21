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
	app.Flags = conf.BuildFlags
	app.Commands = []cli.Command{
		{
			Name:     "build",
			Usage:    "start build and deploy to serverless PaaS",
			Flags:    conf.BuildFlags,
			Before:   conf.MustParseConfig,
			Action:   initServer,
			Category: "BUILD",
		},
		{
			Name:     "gateway",
			Usage:    "start gateway",
			Flags:    conf.GatewayFlags,
			Before:   conf.MustParseConfig,
			Action:   gatewayInit,
			Category: "GATEWAY",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var version string
