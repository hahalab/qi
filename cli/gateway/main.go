package main

import (
	"net/http"
	"os"

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
	aliClient := aliyun.NewClient(conf)
	mux := gateway.NewMux(aliClient)

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, mux)
}
