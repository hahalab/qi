package conf

import "github.com/urfave/cli"

const (
	FlagConfigFile = `config`
)

var BuildFlags = []cli.Flag{
	fileFlag,
}

var GatewayFlags = []cli.Flag{
	fileFlag,
}

var fileFlag = cli.StringFlag{
	Name:   FlagConfigFile + ",c",
	Usage:  "specific a config file `PATH`",
	EnvVar: "CONFIG_PATH",
}
