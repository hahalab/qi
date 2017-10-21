package aliyun

type Service struct {
	Description string    `json:"description"`
	LogConfig   LogConfig `json:"logConfig`
	Role        string    `json:"role"`
	ServiceName string    `json:"serviceName`
}

type LogConfig struct {
	Logstore string `json:"logstore`
	Project  string `json:"project`
}
