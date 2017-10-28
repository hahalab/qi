package config

import (
	"testing"
	"os/user"
	"github.com/sirupsen/logrus"
)

func Test_LoadQiConfig(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		logrus.Fatal(err)
	}
	qiConfig := Config{}
	err = LoadQiConfig(usr.HomeDir+"/.ha.conf", &qiConfig)
	if err != nil {
		panic(err)
	}
}

func Test_LoadCodeConfig(t *testing.T) {
	codePath := "/Users/jiangjinyang/workspace/go/src/gin-demo"
	codeConfig := CodeConfig{}
	err := LoadCodeConfig(codePath+"/qi.yml", &codeConfig)
	if err != nil {
		panic(err)
	}
}
