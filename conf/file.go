package conf

import "strings"

type RawRouterConf []RawRouterLine
type RawRouterLine struct {
	Prefix   string `json:"prefix"`
	Service  string `json:"service"`
	Function string `json:"function"`
}

func (r RawRouterConf) Len() int {
	return len(r)
}

func (r RawRouterConf) Less(i, j int) bool {
	lenOfI := len(strings.SplitN(r[i].Prefix, "/", -1))
	lenOfJ := len(strings.SplitN(r[j].Prefix, "/", -1))

	return lenOfI > lenOfJ
}

func (r RawRouterConf) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
