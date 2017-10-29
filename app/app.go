package app

import (
	"github.com/hahalab/qi/config"
	"github.com/hahalab/qi/aliyun/entity"
	"fmt"
	"strconv"
	"github.com/denverdino/aliyungo/ram"
	"encoding/json"
	"strings"
)

type App struct {
	*config.Config
}

func (app *App) description() string {
	return fmt.Sprintf("%s 由 QI 自动创建, 请避免手动更改", app.Name())
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
		LogConfig: entity.LogConfig{
			Project:  app.ProjectName(),
			Logstore: app.StoreName(),
		},
	}
}

func (app *App) Role() ram.Role {
	assumeRolePolicyDocument := entity.AssumeRolePolicyDocument{
		Statement: []entity.AssumeRolePolicyItem{
			{
				Action: "sts:AssumeRole",
				Effect: "Allow",
				Principal: entity.AssumeRolePolicyPrincpal{
					Service: []string{
						"fc.aliyuncs.com",
					},
				},
			},
		},
		Version: "1",
	}
	assumeRolePolicyDocumentJson, _ := json.Marshal(assumeRolePolicyDocument)
	return ram.Role{
		RoleName:                 app.Name(),
		AssumeRolePolicyDocument: string(assumeRolePolicyDocumentJson),
		Description:              "role is created by QI, 如需调整访问权限请保证角色 Arn 不变, 否则会重新创建",
	}
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
		GroupName:   strings.Replace(fmt.Sprintf("%s-apigroup", app.Name()), "-", "_", -1),
		Description: app.description(),
	}
}

func (app *App) APIGateways(groupId string, roleArn string) ([]entity.APIGateway) {
	return []entity.APIGateway{app.generateAPIGateway(groupId, "GET", roleArn)}
}

func (app *App) generateAPIGateway(groupId string, method string, roleArn string) (entity.APIGateway) {
	return entity.APIGateway{
		RegionId:    app.RegionId(),
		GroupId:     groupId,
		ApiName:     strings.Replace(fmt.Sprintf("%s-%s", app.Name(), method), "-", "_", -1),
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
				RoleArn:             roleArn,
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
