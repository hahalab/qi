package aliyun

type Config struct {
	AccessKeyID     string `validate:"required"`
	AccessKeySecret string `validate:"required"`
	AccountID       string `validate:"required"`
	//cn-beijing.fc.aliyuncs.com
	FcEndPoint string `validate:"required"`
	//oss-cn-beijing.aliyuncs.com
	OssEndPoint   string `validate:"required"`
	OssBucketName string `validate:"required"`
}
