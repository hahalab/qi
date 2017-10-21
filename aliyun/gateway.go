package aliyun

type createApiGateWayReq struct {
	//bcc88ca511be4ad49f568268d67c1c00
	GroupId string
	ApiName string
	//"PUBLIC"
	Visibility  string
	Description string
	//"ANONYMOUS"
	AuthType             string
	RequestConfig        string
	ServiceConfig        string
	RequestParamters     string
	ServiceParameters    string
	ServiceParametersMap string
	ResultType           string
	ResultSample         string
	FailResultSample     string
	ErrorCodeSamples     string
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

type FunctionComputeConfig struct {
	FcRegionId          string
	ServiceName         string
	RoleArn             string
	FunctionName        string
	ContentTypeCatagory string
	ContentTypeValue    string
}

type RequestConfig struct {
	RequestProtocol     string
	RequestHttpMethod   string
	RequestPath         string
	BodyFormat          string
	RequestMode         string
	PostBodyDescription string
}
type RequestParamters struct {
	Required          string
	ParameterType     string
	ApiParameterName  string
	DocShow           string
	Location          string
	DocOrder          int
	ParameterCatalog  string
	ParameterLocation ParameterLocation
	IsHide            bool `json:"isHide"`
}

type ParameterLocation struct {
	Name        string `json:"name"`
	OrderNumber int    `json:"orderNumber"`
}

type ServiceParameters struct {
	ServiceParameterName string
	Location             string
	Type                 string
	ParameterCatalog     string
}

type ServiceParametersMap struct {
	ServiceParameterName string
	RequestParameterName string
}

type createApiGateWayResp struct {
	RequestId string
	ApiId     string
}
