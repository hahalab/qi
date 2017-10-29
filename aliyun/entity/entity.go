package entity

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

type APIGroup struct {
	GroupName   string
	Description string
}

type APIGroupAttribute struct {
	GroupId       string
	GroupName     string
	SubDomain     string
	Description   string
	CreatedTime   string
	ModifiedTime  string
	RegionId      string
	TrafficLimit  int
	BillingStatus string
	IllegalStatus string
}

type APIGateway struct {
	RegionId             string
	GroupId              string
	ApiName              string
	Visibility           string
	Description          string
	AuthType             string
	RequestConfig        RequestConfig
	ServiceConfig        ServiceConfig
	RequestParamters     []string // WTF?
	ServiceParameters    []string
	ServiceParametersMap []string
	ResultType           string
	ResultSample         string
	FailResultSample     string
	ErrorCodeSamples     []string
}

type APISummary struct {
	RegionId     string
	GroupId      string
	ApiId        string
	GroupName    string
	ApiName      string
	Visibility   string
	Description  string
	CreatedTime  string
	ModifiedTime string
}

type ServiceConfig struct {
	ServiceProtocol       string
	ServiceHttpMethod     string
	ServiceAddress        string
	ServiceTimeout        string
	ServicePath           string
	Mock                  string
	MockResult            string
	ServiceVpcEnable      string
	VpcConfig             struct{}
	FunctionComputeConfig FunctionComputeConfig
}

type RequestConfig struct {
	RequestProtocol     string
	RequestHttpMethod   string
	RequestPath         string
	BodyFormat          string
	RequestMode         string
	PostBodyDescription string
}

type FunctionComputeConfig struct {
	FcRegionId          string
	ServiceName         string
	RoleArn             string
	FunctionName        string
	ContentTypeCatagory string
	ContentTypeValue    string
}

// Policy
type AssumeRolePolicyDocument struct {
	Statement []AssumeRolePolicyItem
	Version   string
}

type AssumeRolePolicyItem struct {
	Action    string
	Effect    string
	Principal AssumeRolePolicyPrincpal
}

type AssumeRolePolicyPrincpal struct {
	Service []string
}
