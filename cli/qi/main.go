package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/hahalab/qi/conf"
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
	app.Flags = conf.BuildFlags
	app.Commands = []cli.Command{
		{
			Name:     "build",
			Usage:    "only build code to zip package",
			Flags:    conf.BuildFlags,
			Action:   onlyBuild,
			Category: "BUILD",
		}, {
			Name:     "deploy",
			Usage:    "only create/update function",
			Flags:    conf.BuildFlags,
			Before:   conf.MustParseUpConfig,
			Action:   onlyDeploy,
			Category: "BUILD",
		}, {
			Name:     "gateway",
			Usage:    "start gateway",
			Flags:    conf.GatewayFlags,
			Before:   conf.MustParseGWConfig,
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
