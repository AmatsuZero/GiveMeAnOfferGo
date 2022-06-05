package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const constPrefix = "Contents of (__DATA"
const queryClassList = "__objc_classlist"
const queryClassRefs = "__objc_classrefs"
const querySuperRefs = "__objc_superrefs"
const querySelRefs = "__objc_selrefs"

// 从 ipa 文件生成 Mach-O 文件
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

func checkIsValid(macho string) error {
	valid := false
	err := readFileByLine(macho, func(s string, b *bool) {
		idx := strings.Index(s, constPrefix)
		if idx != -1 {
			valid = true
			*b = true
		}
	})

	if err != nil {
		return err
	} else if !valid {
		return fmt.Errorf("otool文件格式有误")
	}

	return nil
}

/// 获取所有方法集合 { className:{ address: methodName } }
func allSelRefsFromContent(macho string, c *cli.Context) (map[string]map[string]string, error) {
	allSelResultsBegin := false
	canAddName := false
	canAddMethods := false
	className := ""

	methodDict := make(map[string]string)
	allSelResults := make(map[string]map[string]string, 0)

	err := readFileByLine(macho, func(s string, b *bool) {
		if strings.Contains(s, constPrefix) && strings.Contains(s, queryClassList) {
			allSelResultsBegin = true
			return
		} else if allSelResultsBegin && strings.Contains(s, constPrefix) {
			allSelResultsBegin = false
			*b = true
		}

		if !allSelResultsBegin {
			return
		}

		if strings.Contains(s, "data") {
			if len(methodDict) > 0 {
				allSelResults[className] = methodDict
				methodDict = make(map[string]string)
			}
			// data之后第一个的name，是类名
			canAddName = true
			canAddMethods = false
			return
		}

		if canAddName && strings.Contains(s, "name") {
			components := strings.Split(s, " ")
			className = components[len(components)-1]
			return
		}

		if strings.Contains(s, "methods") || strings.Contains(s, "Methods") {
			// method之后的name是方法名，和方法地址
			canAddName = false
			canAddMethods = true
			return
		}

		if canAddMethods && strings.Contains(s, "name") {
			components := strings.Split(s, " ")
			if len(components) > 2 {
				methodAddress := components[len(components)-2]
				methodName := components[len(components)-1]
				methodDict[methodName] = methodAddress
			}
		}
	})

	return allSelResults, err
}

// 获取已使用的方法集合
func selRefsFromContent(macho string, c *cli.Context) (map[string]string, error) {

	selRefsResults := make(map[string]string)
	selRefsBegin := false

	err := readFileByLine(macho, func(s string, b *bool) {
		if strings.Contains(s, constPrefix) && strings.Contains(s, querySelRefs) {
			selRefsBegin = true
			return
		} else if selRefsBegin && strings.Contains(s, constPrefix) {
			selRefsBegin = false
			*b = true
		}

		if !selRefsBegin {
			return
		}

		components := strings.Split(s, " ")
		if len(components) > 2 {
			methodAddress := components[len(components)-2]
			methodName := components[len(components)-1]
			selRefsResults[methodName] = methodAddress
		}
	})

	if c.String("verbose") == "1" {
		fmt.Printf("\n\n__objc_selrefs总结如下，共有%d个\n%v\n：", len(selRefsResults), selRefsResults)
	}

	return selRefsResults, err
}

// 查找多余的方法
func analyzeUsedMethods(macho string, c *cli.Context) string {
	var group sync.WaitGroup
	var methodsListDic map[string]map[string]string
	var selRefsDic map[string]string

	group.Add(1)
	go func() {
		dict, err := allSelRefsFromContent(macho, c)
		if err != nil {
			fmt.Println(err)
		}
		methodsListDic = dict
		group.Done()
	}()

	group.Add(1)
	go func() {
		dict, err := selRefsFromContent(macho, c)
		if err != nil {
			fmt.Println(err)
		}
		selRefsDic = dict
		group.Done()
	}()

	group.Wait()

	// 遍历selRefs移除methodsListDic，剩下的就是未使用的
	for methodAddress, _ := range selRefsDic {
		for _, methodDic := range methodsListDic {
			delete(methodDic, methodAddress)
		}
	}

	// 遍历移除空的元素
	result := make(map[string]map[string]string)
	for classNameStr, methodDic := range methodsListDic {
		if len(methodDic) > 0 {
			result[classNameStr] = methodDic
		}
	}

	if c.String("verbose") == "1" {
		fmt.Printf("多余的方法如下%v", result)
	}

	ret := strings.Builder{}
	ret.WriteString("方法地址\t方法名称\r\n\r\n")
	for className, methodDic := range result {
		ret.WriteString(fmt.Sprintf("%v\t\r\n", className))
		for methodAddress, methodName := range methodDic {
			ret.WriteString(fmt.Sprintf("\t\t\t\t\t%v\t%v\r\n", methodAddress, methodName))
		}
	}

	return ret.String()
}

// 查找多余的类
func analyzeUnusedClass(macho string, c *cli.Context) string {
	// 所有classList类和类名字
	var classListDic map[string]string
	// 所有引用的类
	var classRefs []string

	var group sync.WaitGroup

	group.Add(1)
	go func() {
		dict, err := classListFromContent(macho, c)
		if err != nil {
			fmt.Println(err)
		}
		classListDic = dict
		group.Done()
	}()

	group.Add(1)
	go func() {
		dict, err := classRefsFromContent(macho, c)
		if err != nil {
			fmt.Println(err)
		}
		classRefs = dict
		group.Done()
	}()

	group.Wait()

	// 先把类和父类数组做去重
	classRefs = removeDuplicateStr(classRefs)

	// 所有在refsSet中的都是已使用的，遍历classList，移除refsSet中涉及的类
	// 余下的就是多余的类
	// 移除系统类，比如SceneDelegate，或者Storyboard中的类
	for _, addr := range classRefs {
		delete(classListDic, addr)
	}

	if c.String("verbose") == "1" {
		fmt.Printf("多余的方法如下%v", classListDic)
	}

	ret := strings.Builder{}
	ret.WriteString("文件地址\t文件名称\r\n\r\n")
	for addr, name := range classListDic {
		ret.WriteString(fmt.Sprintf("%v\t%v\r\n", addr, name))
	}
	ret.WriteString(fmt.Sprintf("\r\n总计: %d个\r\n", len(classListDic)))

	return ret.String()
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// 获取 classrefs
func classRefsFromContent(macho string, c *cli.Context) (classRefsResults []string, err error) {
	classRefBegin := false

	err = readFileByLine(macho, func(s string, b *bool) {
		if strings.Contains(s, constPrefix) && strings.Contains(s, queryClassRefs) {
			classRefBegin = true
			return
		} else if classRefBegin && strings.Contains(s, constPrefix) {
			classRefBegin = false
			*b = true
		}

		if classRefBegin && strings.Contains(s, "000000010") {
			components := strings.Split(s, " ")
			addr := components[len(components)-1]
			if strings.HasPrefix(addr, "0x100") {
				classRefsResults = append(classRefsResults, addr)
			}
		}
	})

	if c.String("verbose") == "1" {
		fmt.Printf("\n\n__objc_refs总结如下，共有%d个\n%v\n：", len(classRefsResults), classRefsResults)
	}
	return
}

// 获取 superrefs
func superRefsFromContent(macho string, c *cli.Context) (classSuperRefsResults []string, err error) {
	classSuperRefsBegin := false

	err = readFileByLine(macho, func(s string, b *bool) {
		if strings.Contains(s, constPrefix) && strings.Contains(s, querySuperRefs) {
			classSuperRefsBegin = true
			return
		} else if classSuperRefsBegin && strings.Contains(s, constPrefix) {
			classSuperRefsBegin = false
			*b = true
		}

		if classSuperRefsBegin && strings.Contains(s, "000000010") {
			components := strings.Split(s, " ")
			addr := components[len(components)-1]
			if strings.HasPrefix(addr, "0x100") {
				classSuperRefsResults = append(classSuperRefsResults, addr)
			}
		}
	})

	if c.String("verbose") == "1" {
		fmt.Printf("\n\n__objc_superrefs总结如下，共有%d个\n%v\n：", len(classSuperRefsResults), classSuperRefsResults)
	}
	return
}

// 获取classList的类
func classListFromContent(macho string, c *cli.Context) (map[string]string, error) {
	canAddName := false
	classListBegin := false
	addressStr := ""
	classListResults := make(map[string]string)

	err := readFileByLine(macho, func(s string, b *bool) {
		if strings.Contains(s, constPrefix) && strings.Contains(s, queryClassList) {
			classListBegin = true
			return
		} else if classListBegin && strings.Contains(s, constPrefix) {
			classListBegin = false
			*b = true
		}

		if !classListBegin {
			return
		}

		if strings.Contains(s, "000000010") {
			components := strings.Split(s, " ")
			addressStr = components[len(components)-1]
			canAddName = true
		} else {
			if canAddName && strings.Contains(s, "name") {
				components := strings.Split(s, " ")
				className := components[len(components)-1]
				classListResults[className] = addressStr
				addressStr = ""
				canAddName = false
			}
		}
	})

	if c.String("verbose") == "1" {
		fmt.Printf("__objc_classlist总结如下，共有%d个\n%v\n：", len(classListResults), classListResults)
	}

	return classListResults, err
}

func Unused(c *cli.Context) error {
	if len(os.Args) == 0 || len(os.Args[0]) == 0 {
		return fmt.Errorf("缺少文件路径")
	}
	macho := os.Args[2]

	isTmp := false
	if path.Ext(macho) == ".ipa" {
		path, err := generateMachOFileFromIpa(macho)
		if err != nil {
			return err
		}
		isTmp = true
		macho = path
	}

	err := checkIsValid(macho)
	if err != nil {
		return err
	}

	ret := ""
	if c.String("func") == "1" {
		ret = analyzeUsedMethods(macho, c)
	} else {
		ret = analyzeUnusedClass(macho, c)
	}

	fmt.Println(ret)

	if isTmp { // 移除临时文件
		os.Remove(macho)
	}
	return nil
}
