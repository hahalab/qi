package gateway

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
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

func NewMux(driver driver) http.HandlerFunc {
	return http.HandlerFunc(mux{driver: driver})
}

type Lambda struct {
	PathRegexp *regexp.Regexp
	Path       string
	Name       string
	Service    string
}

func (m mux) Dispatch(w http.ResponseWriter, r *http.Request) {
	lambda, err := m.findLambda(r.RequestURI)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("router compile failed"))
		return
	}

	if lambda == nil {
		w.WriteHeader(404)
		w.Write([]byte("lambda not found"))
		return
	}

	body, _ := json.Marshal(Event{
		Method:  r.Method,
		Headers: r.Header,
		Path:    r.RequestURI,
	})
	respBody, err := m.InvokeFunction(lambda.Service, lambda.Name, body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	resp := Resp{}
	err = json.Unmarshal(respBody, &resp)
	if lambda == nil {
		w.WriteHeader(502)
		w.Write([]byte("unmarshal resp from function failed"))
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

	var conf rawRouterConf

	err = json.NewDecoder(c).Decode(&conf)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	sort.Sort(conf)

	for _, v := range conf {
		l := &Lambda{
			PathRegexp: regexp.MustCompile(strings.Replace(v.URL, "?", "[^/]", -1)),
			Path:       v.URL,
			Name:       v.Name,
			Service:    v.Service,
		}

		if l.PathRegexp.MatchString(url) {
			return l, nil
		}
	}

	return nil, nil
}

type rawRouterConf []rawRouterLine
type rawRouterLine struct {
	URL     string `json:"url"`
	Name    string `json:"name"`
	Service string `json:"service"`
}

func (r rawRouterConf) Len() int {
	return len(r)
}

func (r rawRouterConf) Less(i, j int) bool {
	lenOfI := len(strings.SplitN(r[i].URL, "/", -1))
	lenOfJ := len(strings.SplitN(r[j].URL, "/", -1))

	return lenOfI < lenOfJ
}

func (r rawRouterConf) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type Event struct {
	Method  string
	Path    string
	Headers map[string][]string
	Data    []byte
}

type Resp struct {
	Code    int
	Headers map[string][]string
	Data    []byte
}
