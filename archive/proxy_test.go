package archive

import (
	"testing"
	"archive/zip"
	"bytes"
	"fmt"
)

func Test_ProxyBinData(t *testing.T) {
	input, err := codeZipBytes()
	if err != nil {
		panic(err)
	}
	r, err := zip.NewReader(bytes.NewReader(input), int64(len(input)))
	if err != nil {
		panic(err)
	}
	for _, file := range r.File {
		fmt.Printf("file: %s, size: %i", file.Name, file.UncompressedSize64)
	}
}
