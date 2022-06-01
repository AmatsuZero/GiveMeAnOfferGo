package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const constPrefix = "Contents of (__DATA"
const queryClassList = "__objc_classlist"
const queryClassRefs = "__objc_classrefs"
const querySuperRefs = "__objc_superrefs"
const querySelRefs = "__objc_selrefs"

// 从 ipa 文件生成 mach-O 文件
func generateMachOFileFromIpa(ipaPath string) (string, error) {
	_, fileName := path.Split(ipaPath)

	// 创建临时目录
	dir, err := os.MkdirTemp("", fileName)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	// 重命名为 zip 文件
	zipFilePath := filepath.Join(dir, fileName+".zip")
	err = os.Rename(ipaPath, zipFilePath)
	if err != nil {
		return "", err
	}

	// 最后再改回去
	defer os.Rename(zipFilePath, ipaPath)

	// 解压
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// 找到 Mach-O 文件
	appFolderDirName := ""
	var machoFile *zip.File

	prefix := "Payload/"
	for _, f := range reader.File {
		if len(appFolderDirName) == 0 { // 找到 mach-o 文件名
			if f.FileInfo().IsDir() && strings.HasPrefix(f.Name, prefix) {
				idx := strings.Index(f.Name, prefix)
				name := strings.TrimPrefix(f.Name, prefix)
				if len(name) == 0 {
					continue
				}
				idx = strings.Index(name, ".app")
				if idx != -1 {
					appFolderDirName = name[:idx]
				}
			}
		} else if !f.FileInfo().IsDir() { // 找到对应的文件
			_, file := path.Split(f.Name)
			if file == appFolderDirName {
				machoFile = f
				break
			}
		}
	}

	if machoFile == nil {
		return "", fmt.Errorf("没有找到 mach-o 文件 %v", appFolderDirName)
	}

	zippedFile, err := machoFile.Open()
	if err != nil {
		return "", err
	}
	defer zippedFile.Close()

	dest := filepath.Join(dir, appFolderDirName)
	destinationFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, machoFile.Mode())
	if err != nil {
		return "", err
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return "", err
	}

	cmd := exec.Command("/usr/bin/otool", "-arch", "arm64", "-ov", dest)
	var stderr bytes.Buffer
	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return "", fmt.Errorf(fmt.Sprint(err) + ": " + stderr.String())
	}

	// 生成 txt
	folder, _ := filepath.Split(ipaPath)
	txtPath := filepath.Join(folder, "tmp.txt")
	err = ioutil.WriteFile(txtPath, out.Bytes(), 0755)
	if err != nil {
		return "", err
	}

	return txtPath, nil
}

func checkIsValid(linkMapPath string) error {
	file, err := os.Open(linkMapPath)

	if err != nil {
		return err
	}

	defer file.Close()

	r := bufio.NewReader(file)
	line, isPrefix, err := r.ReadLine()
	for err == nil && !isPrefix {
		s := string(line)
		idx := strings.Index(s, constPrefix)
		if idx != -1 {
			return nil
		}
		line, isPrefix, err = r.ReadLine()
	}

	return fmt.Errorf("otool文件格式有误")
}

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			if len(os.Args) == 0 || len(os.Args[0]) == 0 {
				return fmt.Errorf("缺少文件路径")
			}
			linkMap := os.Args[1]

			isTmp := false
			if path.Ext(linkMap) == ".ipa" {
				path, err := generateMachOFileFromIpa(linkMap)
				if err != nil {
					return err
				}
				isTmp = true
				linkMap = path
			}

			err := checkIsValid(linkMap)
			if err != nil {
				return err
			}

			if isTmp { // 移除临时文件
				os.Remove(linkMap)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
