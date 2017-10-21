package archive

import (
	"archive/zip"
	"path/filepath"
	"os"
	"fmt"
	"strings"
	"io"
)


// 把目录打包成 cwd()/code.zip
func Build(dir string) error {

	output, err := os.Create("code.zip")
	if err != nil {
		return err
	}
	defer output.Close()

	tw := zip.NewWriter(output)

	// Write files to tw
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		zipPath := strings.TrimPrefix(path, dir)
		if len(zipPath) > 1 {
			zipPath = zipPath[1:]
		}
		fmt.Printf("Add file %s to %s\n", path, zipPath)
		header, err := zip.FileInfoHeader(info)

		header.Name = zipPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		if info.Mode()&os.ModeSymlink != 0 {
			header.SetMode(0)
		}

		//fmt.Printf("create heade %+v\n", header)
		writer, err := tw.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
		return nil
	})

	if err := tw.Close(); err != nil {
		return err
	}

	return err
}
