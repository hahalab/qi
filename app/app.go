package app

import (
	"github.com/hahalab/qi/config"
	"github.com/hahalab/qi/aliyun/entity"
	"fmt"
	"strconv"
)

type App struct {
	*config.Config
}

func (app *App) description() string {
	return fmt.Sprintf("%s is deployed by QI, don't manually modified", app.Name())
}

func NewApp(c *config.Config) (*App) {
	return &App{Config: c}
}

func (app *App) Name() string {
	return app.CodeConfig.Name
}

func (app *App) RequestProtocol() string {
	return "HTTP,HTTPS"
}

func (app *App) FunctionTimeout() int {
	return app.CodeConfig.Timeout
}

func (app *App) GatewayTimeout() int {
	return 1000
}

func (app *App) RegionId() string {
	return app.AliyunConfig.RegionId
}

func (app *App) ProjectName() string {
	return fmt.Sprintf("%s-project", app.Name())
}

func (app *App) StoreName() string {
	return fmt.Sprintf("%s-logstore", app.Name())
}

func (app *App) Service() (entity.Service) {
	serviceName := fmt.Sprintf("%s-service", app.Name())
	return entity.Service{
		ServiceName: serviceName,
		Description: app.description(),
		Role:        "",
		LogConfig: entity.LogConfig{
			Project:  app.ProjectName(),
			Logstore: app.StoreName(),
		},
	}
}

func (app *App) Role() string {
	return ""
}

func (app *App) Function() (entity.Function) {
	return entity.Function{
		FunctionName: app.Name(),
		Description:  app.description(),
		MemorySize:   app.CodeConfig.MemorySize,
		Timeout:      app.FunctionTimeout(),
		Handler:      "index.handler",
		Runtime:      "python2.7",
		Code:         entity.Code{},
	}
}

func (app *App) APIGroup() (entity.APIGroup) {
	return entity.APIGroup{
		GroupName:   fmt.Sprintf("%s-apigroup", app.Name()),
		Description: app.description(),
	}
}

func (app *App) APIGateways(groupId string) ([]entity.APIGateway) {
	return []entity.APIGateway{app.generateAPIGateway(groupId, "GET")}
}

func (app *App) generateAPIGateway(groupId string, method string) (entity.APIGateway) {
	return entity.APIGateway{
		RegionId:    app.RegionId(),
		GroupId:     groupId,
		ApiName:     fmt.Sprintf("%s-%s", app.Name(), method),
		Visibility:  "PRIVATE",
		Description: app.description(),
		AuthType:    "ANONYMOUS",
		RequestConfig: entity.RequestConfig{
			RequestProtocol:     app.RequestProtocol(),
			RequestHttpMethod:   method,
			RequestPath:         "/",
			BodyFormat:          "",
			PostBodyDescription: "",
			RequestMode:         "PASSTHROUGH",
		},
		ServiceConfig: entity.ServiceConfig{
			ServiceProtocol:   "FunctionCompute",
			ServiceHttpMethod: method,
			ServiceAddress:    "",
			ServiceTimeout:    strconv.Itoa(app.GatewayTimeout()),
			ServicePath:       "/",
			Mock:              "FALSE",
			MockResult:        "",
			ServiceVpcEnable:  "FALSE",
			VpcConfig:         struct{}{},
			FunctionComputeConfig: entity.FunctionComputeConfig{
				FcRegionId:          app.RegionId(),
				ServiceName:         app.Service().ServiceName,
				FunctionName:        app.Function().FunctionName,
				RoleArn:             app.Role(),
				ContentTypeCatagory: "CLIENT",
				ContentTypeValue:    ""},
		},
		RequestParamters:     nil,
		ServiceParameters:    nil,
		ServiceParametersMap: nil,
		ResultType:           "PASSTHROUGH",
		ResultSample:         "透传 App 响应",
		FailResultSample:     "透传 App 的失败页面",
		ErrorCodeSamples:     nil,
	}
}
