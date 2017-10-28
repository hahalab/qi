package gateway

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/hahalab/qi/conf"
)

type oss interface {
	GetObject(path string) (io.ReadCloser, error)
}

type lambda interface {
	InvokeFunction(service, name string, event []byte) ([]byte, error)
}

type driver interface {
	oss
	lambda
}

type mux struct {
	path string
	driver
}

func NewMux(driver driver, confPath string) http.HandlerFunc {
	return http.HandlerFunc(mux{driver: driver, path: confPath}.Dispatch)
}

type Lambda struct {
	Path    string
	Name    string
	Service string
}

func (m mux) Dispatch(w http.ResponseWriter, r *http.Request) {
	lambda, err := m.findLambda(r.RequestURI)
	if err != nil {
		w.WriteHeader(500)
		logrus.Error("router compile failed")
		return
	}

	if lambda == nil {
		w.WriteHeader(404)
		logrus.Infof("[%s]%s lambda not found", r.Method, r.RequestURI)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		logrus.Debugf("read buff from req body err:%v", err)
		return
	}

	logrus.Debugf("get request: path(%s), found lambda prefix(%s)", r.RequestURI, lambda.Path)

	body, _ := json.Marshal(Event{
		Method:  r.Method,
		Headers: r.Header,
		Path:    r.RequestURI[len(lambda.Path):] + "/",
		Data:    reqBody,
	})

	logrus.Debugf("req body %s, lambda: %v", body, lambda)
	respBody, err := m.InvokeFunction(lambda.Service, lambda.Name, body)
	if err != nil {
		w.WriteHeader(500)
		logrus.Error(err)
		return
	}

	resp := Resp{}
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		w.WriteHeader(502)
		logrus.Infof("%s", respBody)
		logrus.Errorf("unmarshal resp from function failed: %v", err)
		return
	}

	w.WriteHeader(resp.Code)
	for k, v := range resp.Headers {
		for _, h := range v {
			w.Header().Set(k, h)
		}
	}
	w.Write(resp.Data)
}

func (m mux) findLambda(url string) (*Lambda, error) {
	c, err := m.GetObject(m.path)
	if err != nil {
		return nil, err
	}

	var conf conf.RawRouterConf
	err = json.NewDecoder(c).Decode(&conf)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	sort.Sort(conf)

	logrus.Debug("get router conf from oss success")
	for _, v := range conf {
		logrus.Debugf("prefix: %s, service name: %s, func name: %s", v.Prefix, v.Service, v.Function)
	}

	for _, v := range conf {
		if v.Prefix != "" && regexp.MustCompile("^"+v.Prefix).MatchString(url) {
			return &Lambda{
				Path:    v.Prefix,
				Name:    v.Function,
				Service: v.Service,
			}, nil
		}
	}

	return nil, nil
}

type Event struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Data    []byte              `json:"data"`
}

type Resp struct {
	Code    int                 `json:"code"`
	Headers map[string][]string `json:"headers"`
	Data    []byte              `json:"data"`
}
