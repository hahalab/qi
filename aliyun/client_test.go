package aliyun

import (
	"fmt"
	"testing"
)

func Test_CreateService(t *testing.T) {
	cli, err := NewClient(&Config{
		AccessKeyID:     "LTAIII8mgWu95PjV",
		AccessKeySecret: "xxxxxxxxxxxxxxx",
		FcEndPoint:      "cn-beijing.fc.aliyuncs.com",
		AccountID:       "xxxxxxxxxxxxxxx",
		OssBucketName:   "oss-cn-beijing.aliyuncs.com",
	})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(cli.CreateService(Service{
		Description: "testapi",
		Role:        "acs:ram::1759916402662922:role/fc-logs",
		ServiceName: "theapitest",
		LogConfig: LogConfig{
			Logstore: "fc-store",
			Project:  "store",
		},
	}))
}
