package build

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/hahalab/qi/aliyun"
	"github.com/hahalab/qi/archive"
	"github.com/hahalab/qi/conf"
)

type Builder struct {
	*aliyun.Client
}

func NewBuilder(client *aliyun.Client) (Builder, error) {
	return Builder{client}, nil
}

func (b Builder) Build(path string, hintMessage chan string) (err error) {
	err = archive.Build(path, hintMessage)
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}

func (b Builder) Prepare(serviceName, role string, hintMessage chan string) error {
	return nil
}

func (b Builder) Deploy(serviceName, role string, hintMessage chan string) error {
	hintMessage <- "Deploying"

	file, err := ioutil.ReadFile("code.zip")
	if err != nil {
		return err
	}

	fileEncoded := base64.StdEncoding.EncodeToString(file)

	err = b.CreateLogStore("qilog", "fclog")
	if err != nil {
		return err
	}
	err = b.CreateService(aliyun.Service{
		ServiceName: serviceName,
		Description: "s",
		Role:        role,
		LogConfig: aliyun.LogConfig{
			Project:  "qilog",
			Logstore: "fclog",
		},
	})
	if err != nil {
		return err
	}

	err = b.CreateFunction(serviceName, aliyun.Function{
		FunctionName: serviceName + "-func",
		Description:  "f",
		MemorySize:   64 * 4,
		Timeout:      10,
		Handler:      "index.handler",
		Runtime:      "python2.7",
		Code: aliyun.Code{
			fileEncoded,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) Qi(m chan string) error {
	cfg := conf.GetUPConf()

	if err := b.Build(cfg.CodePath, m); err != nil {
		return err
	}

	if err := b.Prepare(cfg.Name, cfg.Role, m); err != nil {
		return err
	}

	if err := b.Deploy(cfg.Name, cfg.Role, m); err != nil {
		return err
	}
	return nil
}
