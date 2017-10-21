package main

import (
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/sirupsen/logrus"
	"github.com/todaychiji/ha/aliyun"
	"github.com/todaychiji/ha/gateway"
)

func main() {
	conf := &aliyun.Config{
		AccessKeyID:     os.Getenv("ALI_ACCESS_KEY"),
		AccessKeySecret: os.Getenv("ALI_ACCESS_SECRET"),
		AccountID:       os.Getenv("ALI_ACCOUNT_ID"),
		Domain:          os.Getenv("ALI_DOMAIN"),
	}
	ok, err := govalidator.ValidateStruct(conf)
	if !ok || err != nil {
		logrus.Fatal(err)
	}

	aliClient := aliyun.NewClient(conf)
	mux := gateway.NewMux(aliClient)

	govalidator.ValidateStruct(conf)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)
}
