package builder

import (
	"github.com/hahalab/qi/aliyun"
	"github.com/hahalab/qi/aliyun/entity"
	"github.com/denverdino/aliyungo/ram"
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
	existsFunction, err := b.client.GetFunction(serviceName, function.FunctionName)

	if existsFunction.FunctionName == "" {
		err = b.client.CreateFunction(serviceName, function)
	} else {
		err = b.client.UpdateFunction(serviceName, function)
	}

	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) EnsureAPIGroup(group entity.APIGroup) (groupAttribute entity.APIGroupAttribute, err error) {
	groupAttribute, err = b.client.GetAPIGroup(group.GroupName)
	if groupAttribute.GroupId != "" {
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

func (b *Builder) EnsureRole(r ram.Role) (role ram.Role, err error) {
	role, err = b.client.GetRole(r.RoleName)
	if role.RoleId != "" && err == nil {
		return
	}
	role, err = b.client.CreateRole(r.RoleName, r.AssumeRolePolicyDocument, r.Description)
	if err != nil {
		return
	}

	// grant policies
	for _, policyName := range []string{
		"AliyunRDSFullAccess",
		"AliyunMongoDBFullAccess",
		"AliyunApiGatewayFullAccess",
		"AliyunLogFullAccess"} {
		_, err := b.client.AttachPolicyToRole(role.RoleName, policyName, "System")
		if err != nil {
			return role, err
		}
	}
	return
}
