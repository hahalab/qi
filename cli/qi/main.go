package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/hahalab/qi/config"
)

func main() {
	app := cli.NewApp()
	app.Name = "qi"
	app.Version = version
	app.Usage = "build and deploy any time any where"
	app.Before = func(context *cli.Context) error {
		logrus.SetOutput(os.Stdout)
		return nil
	}
	app.Action = qi
	app.Flags = config.BuildFlags
	app.Commands = []cli.Command{
		{
			Name:     "build",
			Usage:    "only build code to zip package",
			Flags:    config.BuildFlags,
			Action:   onlyBuild,
			Category: "BUILD",
		}, {
			Name:     "deploy",
			Usage:    "only create/update function",
			Flags:    config.BuildFlags,
			Before:   config.MustParseConfig,
			Action:   onlyDeploy,
			Category: "BUILD",
		}, {
			Name:     "gateway",
			Usage:    "start gateway",
			Flags:    config.GatewayFlags,
			Before:   config.MustParseGWConfig,
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
