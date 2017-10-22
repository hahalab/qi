package main

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/tj/go-spin"
	"github.com/todaychiji/ha/aliyun"
	"github.com/todaychiji/ha/archive"
	"github.com/todaychiji/ha/build"
	"github.com/todaychiji/ha/conf"
)

func qi(c *cli.Context) error {
	message := newMessager()
	message <- "Preparing"

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

	err = build.Qi(message)
	if err != nil {
		logrus.Fatal(err)
	}

	close(message)
	return nil
}

func onlyBuild(c *cli.Context) error {
	message := newMessager()
	defer close(message)
	message <- "Preparing"
	codePath := c.String(conf.FlagCodePath)

	err := archive.Build(codePath, message)
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}

func onlyDeploy(c *cli.Context) error {
	message := newMessager()
	defer close(message)
	message <- "Preparing"

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

	err = build.Deploy(conf.Name, conf.Role, message)
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}

func newMessager() chan string {
	hintMessage := make(chan string, 1)
	go func(m chan string) {
		s := spin.New()

		message := <-m
		for {
			select {
			case message = <-m:
				s.Reset()
				fmt.Printf("\r\033[36m%s\033[m %s ", message, s.Next())
			default:
				fmt.Printf("\r\033[36m%s\033[m %s ", message, s.Next())
				time.Sleep(time.Millisecond * 100)
			}
		}

	}(hintMessage)
	return hintMessage
}
