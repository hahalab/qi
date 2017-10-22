package archive

import "testing"

func Test_Build(t *testing.T) {
	c := make(chan string, 666)
	err := Build("/Users/jiangjinyang/workspace/go/src/github.com/todaychiji/demo", c)
	t.Log(err)
	if err != nil {
		panic(err)
	}
}
