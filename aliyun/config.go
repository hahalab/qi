package aliyun

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	AccountID       string
	//cn-beijing.fc.aliyuncs.com
	FcEndPoint string
	//oss-cn-beijing.aliyuncs.com
	OssEndPoint   string
	OssBucketName string
	//cn-beijing.log.aliyuncs.com
	LogEndPoint string
}
