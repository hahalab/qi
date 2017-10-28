package aliyun

import (
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"github.com/hahalab/qi/config"
	"net"
	"strconv"
	"github.com/denverdino/aliyungo/util"
	"github.com/hahalab/qi/aliyun/entity"
	"encoding/json"
)

type ApiGatewayClient struct {
	conn   *http.Client
	config *config.AliyunConfig
}

var commonQueryParameter map[string]string = map[string]string{
	"Format":           "JSON",
	"Version":          "2016-07-14",
	"SignatureMethod":  "HMAC-SHA1",
	"SignatureVersion": "1.0",
}

func NewApiGatewayClient(config *config.AliyunConfig) (*ApiGatewayClient, error) {
	cli := http.Client{
		Timeout: time.Second * 20,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   20 * time.Second,
				KeepAlive: 120 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          50,
			IdleConnTimeout:       120 * time.Second,
			TLSHandshakeTimeout:   12 * time.Second,
			ExpectContinueTimeout: 2 * time.Second,
			ResponseHeaderTimeout: 12 * time.Second,
		},
	}
	return &ApiGatewayClient{
		conn:   &cli,
		config: config,
	}, nil
}

func endpointURLFromRegionId(regionId string) string {
	return fmt.Sprintf("https://apigateway.%s.aliyuncs.com", regionId)
}

func (client *ApiGatewayClient) get(action string, api interface{}) ([]byte, error) {
	endpoint := endpointURLFromRegionId(client.config.RegionId)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	query := util.ConvertToQueryValues(api)
	// merge common query parameters
	for k, v := range commonQueryParameter {
		query.Add(k, v)
	}
	query.Add("Timestamp", time.Now().UTC().Format(time.RFC3339))
	query.Add("SignatureNonce", strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
	query.Add("AccessKeyId", client.config.AccessKeyID)
	query.Add("Action", action)
	signature := util.CreateSignatureForRequest("GET", &query, client.config.AccessKeySecret+"&")
	query.Add("Signature", signature)
	//build queries
	req.URL.RawQuery = query.Encode()

	resp, err := client.conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("resp status not 200,content:%s", string(content))
	}
	return content, nil
}

type CreateGroupResp struct {
	GroupId     string
	GroupName   string
	SubDomain   string
	Description string
}

// API Gateway

func (client *ApiGatewayClient) GetAPIGroup(groupName string) (groupAttribute *entity.APIGroupAttribute, err error) {
	req := struct {
		GroupName string
	}{
		GroupName: groupName,
	}
	data, err := client.get("DescribeApiGroups", req)
	if err != nil {
		return
	}
	res := struct {
		RequestId          string
		TotalCount         string
		PageSize           string
		PageNumber         string
		ApiGroupAttributes []entity.APIGroupAttribute
	}{}

	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	for _, ga := range res.ApiGroupAttributes {
		if ga.GroupName == groupName {
			groupAttribute = &ga
			break
		}
	}
	return
}

func (client *ApiGatewayClient) CreateAPIGroup(group entity.APIGroup) error {
	result, err := client.get("CreateApiGroup", group)
	fmt.Printf("%s", result)
	if err != nil {
		return err
	}
	return nil
}

func (client *ApiGatewayClient) GetAPIGateway(groupId string, apiName string) (apiSummary *entity.APISummary, err error) {
	req := struct {
		GroupId string
		ApiName string
	}{
		GroupId: groupId,
		ApiName: apiName,
	}
	data, err := client.get("DescribeApis", req)
	if err != nil {
		return
	}
	res := struct {
		RequestId   string
		TotalCount  string
		PageSize    string
		PageNumber  string
		ApiSummarys []entity.APISummary
	}{}

	err = json.Unmarshal(data, res)
	if err != nil {
		return
	}
	for _, api := range res.ApiSummarys {
		if api.GroupId == groupId && api.ApiName == apiName {
			apiSummary = &api
			break
		}
	}
	return
}

func (client *ApiGatewayClient) CreateAPIGateway(api entity.APIGateway) error {
	result, err := client.get("CreateApi", api)
	fmt.Printf("%s", result)
	if err != nil {
		return err
	}
	return nil
}
