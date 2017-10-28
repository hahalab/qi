package main

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/go-playground/validator.v9"

	"github.com/hahalab/qi/aliyun"
	"github.com/hahalab/qi/conf"
	"github.com/hahalab/qi/gateway"
)

func gatewayInit(c *cli.Context) {
	conf := conf.GetGWConf()

	err := validator.New().Struct(conf)
	if err != nil {
		logrus.Fatal(err)
	}

	aliClient, err := aliyun.NewClient(&conf.AliyunConfig)
	if err != nil {
		logrus.Fatal(err)
	}

	mux := gateway.NewMux(aliClient, conf.RouterPath)

	govalidator.ValidateStruct(conf)

	logrus.Infof("gateway listen at :%s", conf.Port)
	http.ListenAndServe(":"+conf.Port, mux)
}
