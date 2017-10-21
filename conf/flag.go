package conf

import "github.com/urfave/cli"

const (
	FlagConfigFile = `config`
)

var Flags = []cli.Flag{
	cli.StringFlag{
		Name:   FlagConfigFile + ",c",
		Value:  "./up.conf",
		Usage:  "specific a config file `PATH`",
		EnvVar: "CONFIG_PATH",
	},
}
