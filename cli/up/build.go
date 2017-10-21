package main

import (
	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/todaychiji/ha/aliyun"
	"github.com/todaychiji/ha/archive"
	"github.com/todaychiji/ha/build"
	"github.com/todaychiji/ha/conf"
)

func up(c *cli.Context) error {
	if err := conf.MustParseUpConfig(c); err != nil {
		return err
	}

	conf := conf.GetUPConf()

	err := validator.New().Struct(conf)
	if err != nil {
		logrus.Fatal(err)
	}

	aliClient, err := aliyun.NewClient(&conf.Config)
	if err != nil {
		logrus.Fatal(err)
	}

	build, err := build.NewBuilder(aliClient)
	if err != nil {
		return err
	}

	return build.Up()
}

func onlyBuild(c *cli.Context) error {
	codePath := c.String(conf.FlagCodePath)
	return archive.Build(codePath)
}

func onlyDeploy(c *cli.Context) error {
	conf := conf.GetUPConf()

	err := validator.New().Struct(conf)
	if err != nil {
		logrus.Fatal(err)
	}

	aliClient, err := aliyun.NewClient(&conf.Config)
	if err != nil {
		logrus.Fatal(err)
	}

	build, err := build.NewBuilder(aliClient)
	if err != nil {
		return err
	}

	return build.Deploy(conf.Name)
}
