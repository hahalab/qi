package archive

import (
	"testing"
	"archive/zip"
	"bytes"
)

func Test_ProxyBinData(t *testing.T) {
	input, err := codeZipBytes()
	if err != nil {
		panic(err)
	}
	_, err = zip.NewReader(bytes.NewReader(input), int64(len(input)))
	if err != nil {
		panic(err)
	}
}
