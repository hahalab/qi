package aliyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hahalab/qi/config"
)

type Client struct {
	conn   *http.Client
	config *config.AliyunConfig
	ossCli *oss.Client
	logCli *sls.Client
	apiCli *ApiGatewayClient
}

func NewClient(config *config.AliyunConfig) (*Client, error) {
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
	ossCli, err := oss.New("http://"+config.OssEndPoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	err = ossCli.CreateBucket(config.OssBucketName)
	if err != nil {
		return nil, err
	}
	logCli := &sls.Client{Endpoint: config.LogEndPoint, AccessKeyID: config.AccessKeyID, AccessKeySecret: config.AccessKeySecret}

	apiCLi, err := NewApiGatewayClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   &cli,
		config: config,
		ossCli: ossCli,
		logCli: logCli,
		apiCli: apiCLi,
	}, nil
}

func (client *Client) get(path string) ([]byte, error) {
	host := fmt.Sprintf("http://%s.%s", client.config.AccountID, client.config.FcEndPoint)
	req, err := http.NewRequest("GET", host+path, nil)
	if err != nil {
		return nil, err
	}
	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set(HTTPHeaderDate, date)
	client.signHeader(req, path)
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

func (client *Client) delete(path string) error {
	host := fmt.Sprintf("http://%s.%s", client.config.AccountID, client.config.FcEndPoint)
	req, err := http.NewRequest("DELETE", host+path, nil)
	if err != nil {
		return err
	}
	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set(HTTPHeaderDate, date)
	client.signHeader(req, path)
	resp, err := client.conn.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("resp status not 200,content:%s", string(content))
	}
	return nil
}

func (client *Client) post(path string, reqBody []byte) ([]byte, error) {
	host := fmt.Sprintf("http://%s.%s", client.config.AccountID, client.config.FcEndPoint)
	req, err := http.NewRequest("POST", host+path, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set(HTTPHeaderDate, date)
	req.Header.Set(HTTPHeaderContentType, "application/json")
	client.signHeader(req, path)
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

func (client *Client) CreateService(service Service) error {
	_, err := client.get(fmt.Sprintf("/2016-08-15/services/%s", service.ServiceName))
	if err == nil {
		return nil
	}
	reqBody, err := json.Marshal(service)
	if err != nil {
		return err
	}
	_, err = client.post("/2016-08-15/services", reqBody)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) CreateFunction(serviceName string, function Function) error {
	_, err := client.get(fmt.Sprintf("/2016-08-15/services/%s/functions/%s", serviceName, function.FunctionName))
	if err == nil {
		err = client.delete(fmt.Sprintf("/2016-08-15/services/%s/functions/%s", serviceName, function.FunctionName))
		if err != nil {
			return err
		}
	}

	reqBody, err := json.Marshal(function)
	if err != nil {
		return err
	}
	_, err = client.post(fmt.Sprintf("/2016-08-15/services/%s/functions", serviceName), reqBody)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) InvokeFunction(serviceName string, functionName string, event []byte) ([]byte, error) {
	content, err := client.post(fmt.Sprintf("/2016-08-15/services/%s/functions/%s/invocations", serviceName, functionName), event)
	if err != nil {
		return nil, err
	}
	return content, err
}

//BucketName为阿里云全局唯一，所以需要加NameSpace避免命名冲突
func (client *Client) CreateBucket(bucketName string) error {
	return client.ossCli.CreateBucket(bucketName)
}

func (client *Client) PutObject(objectKey string, reader io.Reader) error {
	bucket, err := client.ossCli.Bucket(client.config.OssBucketName)
	if err != nil {
		return err
	}
	return bucket.PutObject(objectKey, reader)
}

func (client *Client) GetObject(objectKey string) (io.ReadCloser, error) {
	bucket, err := client.ossCli.Bucket(client.config.OssBucketName)
	if err != nil {
		return nil, err
	}
	return bucket.GetObject(objectKey)
}

func (client *Client) CreateLogStore(ProjectName string, StoreName string) error {
	isExsist, err := client.logCli.CheckProjectExist(ProjectName)
	if err != nil {
		return err
	}
	var p *sls.LogProject
	if isExsist {
		p, err = client.logCli.GetProject(ProjectName)
	} else {
		p, err = client.logCli.CreateProject(ProjectName, "ha log for function compute")
	}
	if err != nil {
		return err
	}
	logExsist, err := p.CheckLogstoreExist(StoreName)
	if err != nil {
		return err
	}
	if !logExsist {
		err = p.CreateLogStore(StoreName, 30, 1)
		if err != nil {
			return err
		}
	}
	return nil
}
