package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
)

type symbol struct {
	//文件
	file string
	// 大小
	size uint64
}

func isValidLinkMap(file string) error {
	containsObjsFileTag := false
	symbolsRange := false
	containsPath := false
	err := readFileByLine(file, func(s string, b *bool) {
		if strings.Contains(s, "# Object files:") {
			containsObjsFileTag = true
		} else if strings.Contains(s, "# Symbols:") {
			symbolsRange = true
		} else if strings.Contains(s, "# Path:") {
			containsPath = true
		}
	})

	if err != nil {
		return err
	}

	if containsObjsFileTag && symbolsRange && containsPath {
		return nil
	}

	return fmt.Errorf("不合法的内容")
}

func ClassSize(c *cli.Context) error {

	return nil
}
