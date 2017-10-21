package aliyun

import (
	"fmt"
	"testing"
)

func Test_CreateService(t *testing.T) {
	cli := NewClient(&Config{
		AccessKeyID:     "LTAIII8mgWu95PjV",
		AccessKeySecret: "xxxxxxxxxxxxxxx",
		Domain:          "cn-beijing.fc.aliyuncs.com",
		AccountID:       "xxxxxxxxxxxxxxx",
	})
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
