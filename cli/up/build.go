package main

import (
	"fmt"

	"github.com/todaychiji/ha/conf"
	"github.com/urfave/cli"
)

func initServer(c *cli.Context) error {
	fmt.Println(c.String(conf.FlagConfigFile))
	return nil
}
