package archive

import (
	"archive/zip"
	"path/filepath"
	"os"
	"os/exec"
	"fmt"
	"strings"
	"io"
	"github.com/todaychiji/ha/config"
	"path"
	"bytes"
)

// 把目录打包成 cwd()/code.zip
func Build(dir string) error {
	c := config.LoadConfig(path.Join(dir, "ha.yml"))

	if err := executeBuild(dir, c); err != nil {
		return err
	}

	output, err := os.Create("code.zip")
	if err != nil {
		return err
	}
	defer output.Close()

	tw := zip.NewWriter(output)
	// Write index.py to zip
	err = injectProxy(tw)
	if err != nil {
		return err
	}
	// Write files to tw
	err = injectDir(dir, tw)

	if err := tw.Close(); err != nil {
		return err
	}

	return err
}

func injectProxy(tw *zip.Writer) error {
	dir := "/Users/jiangjinyang/workspace/go/src/github.com/todaychiji/ha/proxy"
	return injectDir(dir, tw)
}

func injectDir(dir string, tw *zip.Writer) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

	return err
}

func executeBuild(dir string, c config.Config) error {

	if c.Build == "" {
		return nil
	}

	oldPwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err = os.Chdir(dir); err != nil {
		return err
	}

	cmd := exec.Command("sh", "-c", c.Build)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	fmt.Printf("Build result: %q\n", out.String())
	if err != nil {
		return err
	}

	if err = os.Chdir(oldPwd); err != nil {
		return err
	}
	return nil
}
