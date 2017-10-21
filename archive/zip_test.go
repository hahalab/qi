package archive

import "testing"

func Test_Build(t *testing.T) {
	err := Build("/Users/jiangjinyang/workspace/test_fc/code")
	t.Log(err)
	if err != nil {
		panic(err)
	}
}
