package aliyun

import (
	"encoding/json"
	"regexp"
)

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

type deployApiGateWayReq struct {
	GroupId     string
	ApiId       string
	StageName   string
	Description string
}
type deployApiGateWayResp struct {
	RequestId string
}
type createGroupReq struct {
	GroupName   string
	Description string
}
type createGroupResp struct {
	GroupId     string
	GroupName   string
	SubDomain   string
	Description string
}

func (client *Client) createGroup(groupName string) (*createGroupResp, error) {
	var req createGroupReq = createGroupReq{
		GroupName:   groupName,
		Description: "fc gateway",
	}
	var resp createGroupResp
	err := client.commonCli.Invoke("CreateApi", req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (client *Client) createApiGateWay(method string, path string, serviceName string, functionName string, roleArn string, regionId string, groupId string) error {
	reqConfigStr, _ := json.Marshal(RequestConfig{
		RequestProtocol:     "HTTP",
		RequestHttpMethod:   method,
		RequestPath:         path,
		BodyFormat:          "",
		PostBodyDescription: "",
		RequestMode:         "PASSTHROUGH",
	})
	serviceConfigStr, _ := json.Marshal(ServiceConfig{
		ServiceProtocol:   "FunctionCompute",
		ServiceHttpMethod: method,
		ServiceAddress:    "",
		ServiceTimeout:    "5000",
		ServicePath:       "/",
		Mock:              "FALSE",
		MockResult:        "",
		ServiceVpcEnable:  "FALSE",
		VpcConfig:         struct{}{},
		FunctionComputeConfig: FunctionComputeConfig{
			FcRegionId:          regionId,
			ServiceName:         serviceName,
			RoleArn:             roleArn,
			FunctionName:        functionName,
			ContentTypeCatagory: "CLIENT",
			ContentTypeValue:    "",
		},
	})
	r, _ := regexp.Compile(`/\[(\w+)\]`)
	parameters := r.FindAllString(path, -1)
	rps := make([]RequestParamters, 0)
	sps := make([]ServiceParameters, 0)
	spms := make([]ServiceParametersMap, 0)
	for i := range parameters {
		para := parameters[i]
		rp := RequestParamters{
			Required:         "REQUIRED",
			ParameterType:    "String",
			ApiParameterName: para[2 : len(para)-1],
			DocShow:          "PUBLIC",
			Location:         "Path",
			DocOrder:         0,
			ParameterCatalog: "REQUEST",
			IsHide:           false,
			ParameterLocation: ParameterLocation{
				Name:        "Parameter Path",
				OrderNumber: 1,
			},
		}
		sp := ServiceParameters{
			ServiceParameterName: para,
			Location:             "Path",
			Type:                 "String",
			ParameterCatalog:     "REQUEST",
		}
		spm := ServiceParametersMap{
			ServiceParameterName: para,
			RequestParameterName: para,
		}
		rps = append(rps, rp)
		sps = append(sps, sp)
		spms = append(spms, spm)
	}
	rpsStr, _ := json.Marshal(rps)
	spsStr, _ := json.Marshal(sps)
	spmsStr, _ := json.Marshal(spms)
	var req createApiGateWayReq = createApiGateWayReq{
		GroupId:              groupId,
		ApiName:              serviceName + "-" + functionName,
		Visibility:           "PUBLIC",
		Description:          "thisisatest",
		AuthType:             "ANONYMOUS",
		RequestConfig:        string(reqConfigStr),
		ServiceConfig:        string(serviceConfigStr),
		RequestParamters:     string(rpsStr),
		ServiceParameters:    string(spsStr),
		ServiceParametersMap: string(spmsStr),
		ResultType:           "JSON",
		ResultSample:         `{"bodyParam":"testBody"}`,
		FailResultSample:     `{"bodyParam":"testBody"}`,
		ErrorCodeSamples:     "[]",
	}
	var resp createApiGateWayResp
	err := client.commonCli.Invoke("CreateApi", req, &resp)
	if err != nil {
		return err
	}
	var deployReq deployApiGateWayReq = deployApiGateWayReq{
		GroupId:     groupId,
		ApiId:       resp.ApiId,
		StageName:   "RELEASE",
		Description: "test",
	}
	var deployResp deployApiGateWayResp
	err = client.commonCli.Invoke("DeployApi", deployReq, &deployResp)
	if err != nil {
		return err
	}
	return nil
}
