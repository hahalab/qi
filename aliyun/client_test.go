package aliyun

import (
	"os"
	"testing"
)

func newCli() *Client {
	aks := os.Getenv("AccessKeySecret")
	accoundId := os.Getenv("AccountID")
	cli, err := NewClient(&Config{
		AccessKeyID:     "LTAIII8mgWu95PjV",
		AccessKeySecret: aks,
		FcEndPoint:      "cn-beijing.fc.aliyuncs.com",
		AccountID:       accoundId,
		OssBucketName:   "this-test",
		OssEndPoint:     "oss-cn-beijing.aliyuncs.com",
		LogEndPoint:     "cn-beijing.log.aliyuncs.com",
		ApiEndPoint:     "apigateway.cn-beijing.aliyuncs.com",
	})
	if err != nil {
		panic(err)
	}
	return cli
}

func Test_CreateService(t *testing.T) {
	cli := newCli()
	err := cli.CreateService(Service{
		Description: "testapi",
		Role:        "acs:ram::1759916402662922:role/fc-logs",
		ServiceName: "theapitest",
		LogConfig: LogConfig{
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
