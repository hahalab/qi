package aliyun

import (
	"testing"
	"github.com/hahalab/qi/config"
	"os/user"
	"github.com/sirupsen/logrus"
	"github.com/hahalab/qi/aliyun/entity"
	"fmt"
)

func newCli() *Client {
	usr, err := user.Current()
	if err != nil {
		logrus.Fatal(err)
	}
	qiConfig := config.Config{}
	err = config.LoadQiConfig(usr.HomeDir+"/.ha.conf", &qiConfig)

	cli, err := NewClient(&(qiConfig.CommonConf.AliyunConfig))
	if err != nil {
		panic(err)
	}
	return cli
}

func Test_CreateService(t *testing.T) {
	cli := newCli()
	err := cli.CreateService(entity.Service{
		Description: "testapi",
		Role:        "acs:ram::1759916402662922:role/fc-logs",
		ServiceName: "theapitest",
		LogConfig: entity.LogConfig{
			Logstore: "fc-store",
			Project:  "fc-store-test-it",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_CreateLog(t *testing.T) {
	cli := newCli()
	err := cli.CreateLogStore("fc-store-test-it", "fc-store")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_CreateAPIGatewayGroup(t *testing.T) {
	cli := newCli()
	err := cli.CreateAPIGroup(entity.APIGroup{
		GroupName:   "qitest",
		Description: "qi test",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetAPIGatewayGroup(t *testing.T) {
	cli := newCli()
	group, err := cli.GetAPIGroup("qitest")
	t.Logf("%+v", group)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_APIGateway(t *testing.T) {
	cli := newCli()
	err := cli.CreateAPIGateway(entity.APIGateway{
		RegionId:    "cn-shanghai",
		GroupId:     "fb64de791f7a4c708e7a97a2c5e7172d",
		ApiName:     "qitest",
		Visibility:  "PRIVATE",
		Description: "qitest",
		AuthType:    "ANONYMOUS",
		RequestConfig: entity.RequestConfig{
			RequestProtocol:     "HTTP,HTTPS",
			RequestHttpMethod:   "PATCH",
			RequestPath:         "/",
			BodyFormat:          "",
			PostBodyDescription: "",
			RequestMode:         "PASSTHROUGH",
		},
		ServiceConfig: entity.ServiceConfig{
			ServiceProtocol:   "FunctionCompute",
			ServiceHttpMethod: "GET",
			ServiceAddress:    "",
			ServiceTimeout:    "500",
			ServicePath:       "/",
			Mock:              "FALSE",
			MockResult:        "",
			ServiceVpcEnable:  "FALSE",
			VpcConfig:         struct{}{},
			FunctionComputeConfig: entity.FunctionComputeConfig{
				FcRegionId:          "cn-shanghai",
				ServiceName:         "test",
				FunctionName:        "testyaml",
				RoleArn:             "acs:ram::1896697416215058:role/aliyunapigatewayaccessingfcrole",
				ContentTypeCatagory: "CLIENT",
				ContentTypeValue:    ""},
		},
		RequestParamters:     nil,
		ServiceParameters:    nil,
		ServiceParametersMap: nil,
		ResultType:           "PASSTHROUGH",
		ResultSample:         "asd",
		FailResultSample:     "asd",
		ErrorCodeSamples:     nil,
	})
	if err != nil {
		t.Fatal(err)
	}
}
