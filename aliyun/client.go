package aliyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Client struct {
	conn   *http.Client
	config *Config
}

func NewClient(config *Config) *Client {
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

	return &Client{
		conn:   &cli,
		config: config,
	}
}

func (client *Client) Get(path string) ([]byte, error) {
	host := fmt.Sprintf("http://%s.%s", client.config.AccountID, client.config.Domain)
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("resp status not 200")
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (client *Client) Post(path string, reqBody []byte) ([]byte, error) {
	host := fmt.Sprintf("http://%s.%s", client.config.AccountID, client.config.Domain)
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
	fmt.Println(string(content))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("resp status not 200")
	}
	return content, nil
}

func (client *Client) CreateService(service Service) error {
	reqBody, err := json.Marshal(service)
	if err != nil {
		return err
	}
	_, err = client.Post("/2016-08-15/services", reqBody)
	if err != nil {
		return err
	}
	return nil
}
