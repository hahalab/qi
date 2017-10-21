package gateway

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type oss interface {
	GetFile(path string) ([]byte, error)
}

type lambda interface {
	InvokeLambda(service, name string, event Event) (Resp, error)
}

var o oss
var l lambda

type Lambda struct {
	PathRegexp *regexp.Regexp
	Path       string
	Name       string
	Service    string
}

func Dispatch(w http.ResponseWriter, r *http.Request) {
	lambda, err := findLambda(r.RequestURI)
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

	resp, err := l.InvokeLambda(lambda.Service, lambda.Name, Event{})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(resp.Code)
	for k, v := range resp.Headers {
		w.Header().Set(k, v)
	}
	w.Write(resp.Data)
}

func findLambda(url string) (*Lambda, error) {
	c, err := o.GetFile("router.conf")
	if err != nil {
		return nil, err
	}

	var conf rawRouterConf

	err = json.Unmarshal(c, &conf)
	if err != nil {
		return nil, err
	}

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
	Headers map[string]string
	Data    []byte
}

type Resp struct {
	Code    int
	Headers map[string]string
	Data    []byte
}
