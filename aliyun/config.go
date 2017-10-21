package aliyun

type Config struct {
	AccessKeyID     string `validate:"required"`
	AccessKeySecret string `validate:"required"`
	AccountID       string `validate:"required"`
	Domain          string `validate:"required"`
}
