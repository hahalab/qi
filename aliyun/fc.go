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

type APIGroup struct {
	GroupName   string
	Description string
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

//RegionId:cn-shanghai
//GroupId:fb64de791f7a4c708e7a97a2c5e7172d
//ApiName:haha
//Visibility:PRIVATE
//Description:dfg
//AuthType:ANONYMOUS
//RequestConfig:{"RequestProtocol":"HTTP,HTTPS","RequestHttpMethod":"GET","RequestPath":"/","BodyFormat":"","PostBodyDescription":"","RequestMode":"PASSTHROUGH"}
//ServiceConfig:{"ServiceProtocol":"FunctionCompute","ServiceHttpMethod":"GET","ServiceAddress":"","ServiceTimeout":"500","ServicePath":"/","Mock":"FALSE","MockResult":"","ServiceVpcEnable":"FALSE","VpcConfig":{},"FunctionComputeConfig":{"FcRegionId":"cn-shanghai","ServiceName":"test","FunctionName":"testyaml","RoleArn":"acs:ram::1896697416215058:role/aliyunapigatewayaccessingfcrole"},"ContentTypeCatagory":"CLIENT","ContentTypeValue":""}
//RequestParamters:[]
//ServiceParameters:[]
//ServiceParametersMap:[]
//ResultType:PASSTHROUGH
//ResultSample:asd
//FailResultSample:asd
//ErrorCodeSamples:[]
//OpenIdConnectConfig:undefined
//secToken:4KukFEU1WrMNwt79UFAfAC
