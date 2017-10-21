package aliyun

type Service struct {
	Description string    `json:"description"`
	LogConfig   LogConfig `json:"logConfig"`
	Role        string    `json:"role"`
	ServiceName string    `json:"serviceName"`
}

type LogConfig struct {
	Logstore string `json:"logstore"`
	Project  string `json:"project"`
}

type Function struct {
	Description  string `json:"description"`
	FunctionName string `json:"functionName"`
	Handler      string `json:"handler"`
	MemorySize   int    `json:"memorySize"`
	Runtime      string `json:"runtime"`
	Timeout      int    `json:"timeout"`
	Code         Code   `json:"code"`
}

type Code struct {
	ZipFile string `json:"zipFile"`
}
