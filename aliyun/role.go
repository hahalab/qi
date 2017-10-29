package aliyun

import (
	"github.com/hahalab/qi/config"
	"github.com/denverdino/aliyungo/ram"
)

type RoleClient struct {
	client ram.RamClientInterface
	config *config.AliyunConfig
}

func NewRoleClient(config *config.AliyunConfig) (*RoleClient, error) {
	client := ram.NewClient(config.AccessKeyID, config.AccessKeySecret)
	return &RoleClient{
		client: client,
		config: config,
	}, nil
}

func (client *RoleClient) GetRole(roleName string) (role ram.Role, err error) {
	resp, err := client.client.GetRole(ram.RoleQueryRequest{RoleName: roleName,})
	if err != nil {
		return
	}
	role = resp.Role
	return
}

func (client *RoleClient) CreateRole(roleName string, assumeRolePolicyDocument string, description string) (role ram.Role, err error) {
	resp, err := client.client.CreateRole(
		ram.RoleRequest{
			RoleName:                 roleName,
			AssumeRolePolicyDocument: assumeRolePolicyDocument,
			Description:              description,
		},
	)
	if err != nil {
		return
	}
	role = resp.Role
	return
}

func (client *RoleClient) AttachPolicyToRole(roleName string, policyName string, policyType string) (resp ram.RamCommonResponse, err error) {
	resp, err = client.client.AttachPolicyToRole(
		ram.AttachPolicyToRoleRequest{
			PolicyRequest: ram.PolicyRequest{
				PolicyName: policyName,
				PolicyType: ram.Type(policyType),
				//Description:    description,
				//PolicyDocument: policyDocument,
				//SetAsDefault:   setAsDefault,
				//VersionId:      versionId,
			},
			RoleName: roleName,
		},
	)
	return
}
