package conf

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	FlagCodePath = `dir`
	FlagPortPath = `port`
)

var pwd string

func init() {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}
}

var BuildFlags = []cli.Flag{
	fileFlag,
}

var GatewayFlags = []cli.Flag{
	cli.StringFlag{
		Name:   FlagPortPath + ",p",
		Usage:  "specific a `PORT` to listen",
		EnvVar: "PORT",
		Value:  "8080",
	},
}

var fileFlag = cli.StringFlag{
	Name:   FlagCodePath + ",d",
	Usage:  "specific a code `PATH`",
	EnvVar: "CODE_PATH",
	Value:  pwd,
}
