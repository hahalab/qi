package main

import (
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"

	"github.com/todaychiji/ha/aliyun"
	"github.com/todaychiji/ha/gateway"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	conf := &aliyun.Config{
		AccessKeyID:     os.Getenv("ALI_ACCESS_KEY"),
		AccessKeySecret: os.Getenv("ALI_ACCESS_SECRET"),
		AccountID:       os.Getenv("ALI_ACCOUNT_ID"),
		FcEndPoint:      os.Getenv("ALI_FC_ENDPOINT"),
		OssEndPoint:     os.Getenv("ALI_OSS_ENDPOINT"),
		OssBucketName:   os.Getenv("ALI_OSS_BUCKET_NAME"),
		LogEndPoint:     os.Getenv("ALI_LOG_ENDPOINT"),
	}
	err := validator.New().Struct(conf)
	if err != nil {
		logrus.Fatal(err)
	}

	aliClient, err := aliyun.NewClient(conf)
	if err != nil {
		logrus.Fatal(err)
	}
	confPath := os.Getenv("ROUTER_PATH")
	if confPath == "" {
		confPath = "router.conf"
	}
	mux := gateway.NewMux(aliClient, confPath)

	govalidator.ValidateStruct(conf)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)
}
