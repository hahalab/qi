package build

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/todaychiji/ha/aliyun"
	"github.com/todaychiji/ha/archive"
	"github.com/todaychiji/ha/conf"
)

type Builder struct {
	*aliyun.Client
}

func NewBuilder(client *aliyun.Client) (Builder, error) {
	return Builder{client}, nil
}

func (b Builder) Build(path string) (err error) {
	return archive.Build(path)
}

func (b Builder) Deploy(serviceName string) error {
	file, err := ioutil.ReadFile("code.zip")
	if err != nil {
		return err
	}

	fileEncoded := base64.StdEncoding.EncodeToString(file)

	err = b.CreateService(aliyun.Service{
		ServiceName: serviceName,
		Description: "s",
		Role:        "acs:ram::1759916402662922:role/fc-logs",
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

func (b Builder) Up() error {
	cfg := conf.GetUPConf()
	routerReader, err := b.GetObject(cfg.RouterPath)
	if err != nil {
		return err
	}

	//update function
	if err := b.Build(cfg.CodePath); err != nil {
		return err
	}

	if err := b.Deploy(cfg.Name); err != nil {
		return err
	}

	//update gateway routers
	routerRouters := conf.RawRouterConf{}
	if err := json.NewDecoder(routerReader).Decode(&routerRouters); err != nil {
		return err
	}

	var index = -1
	for i := range routerRouters {
		index = i
	}
	if index == -1 {
		routerRouters = append(routerRouters, conf.RawRouterLine{})
		index = len(routerRouters) - 1
	}
	routerRouters[index] = conf.RawRouterLine{
		Prefix:   "/" + cfg.Name,
		Service:  cfg.Name,
		Function: cfg.Name + "-func",
	}

	buf, err := json.Marshal(routerRouters)
	if err != nil {
		return err
	}
	logrus.Debugf("%s", buf)

	err = b.PutObject(cfg.RouterPath, bytes.NewReader(buf))
	if err != nil {
		return err
	}

	return nil
}
