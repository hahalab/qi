package archive

import "testing"

func Test_Build(t *testing.T) {
	err := Build("/Users/jiangjinyang/workspace/go/src/gin-demo")
	t.Log(err)
	if err != nil {
		panic(err)
	}
}
