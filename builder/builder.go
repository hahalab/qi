package builder

import (
	"github.com/hahalab/qi/aliyun"
	"github.com/hahalab/qi/aliyun/entity"
	"io/ioutil"
	"encoding/base64"
)

type Builder struct {
	client *aliyun.Client
}

func NewBuilder(client *aliyun.Client) (Builder, error) {
	return Builder{client}, nil
}

func (b *Builder) EnsureLogStore(projectName string, storeName string) (err error) {
	err = b.client.CreateLogStore(projectName, storeName)
	return
}

func (b *Builder) EnsureService(service entity.Service) (err error) {
	err = b.client.CreateService(service)
	return
}

func (b *Builder) DeployFunction(serviceName string, function entity.Function) (err error) {
	file, err := ioutil.ReadFile("code.zip")
	if err != nil {
		return err
	}

	function.Code.ZipFile = base64.StdEncoding.EncodeToString(file)

	existsFunction, err := b.client.GetFunction(serviceName, function.FunctionName)
	if err != nil || existsFunction == nil {
		err = b.client.CreateFunction(serviceName, function)
	} else {
		err = b.client.UpdateFunction(serviceName, function)
	}

	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) EnsureAPIGroup(group entity.APIGroup) (groupAttribute *entity.APIGroupAttribute, err error) {
	groupAttribute, err = b.client.GetAPIGroup(group.GroupName)
	if groupAttribute != nil && err == nil {
		return
	}
	err = b.client.CreateAPIGroup(group)
	if err != nil {
		return
	}
	groupAttribute, err = b.client.GetAPIGroup(group.GroupName)
	return
}

func (b *Builder) EnsureAPIGateway(groupId string, api entity.APIGateway) (err error) {
	apiSummary, err := b.client.GetAPIGateway(groupId, api.ApiName)
	if apiSummary != nil && err == nil {
		return
	}
	err = b.client.CreateAPIGateway(api)
	return
}
