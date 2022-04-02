package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dir := "./sample"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("ioutil.ReadDir error")
		os.Exit(1)
	}
	var toRename []string
	count := 0
	for _, file := range files {
		if !file.IsDir() {
			//如果是文件，就获取文件名并进行match
			_, err := match(file.Name(), 4)
			if err == nil {
				count++
				toRename = append(toRename, file.Name())
			}
		}
	}
	for _, origFilename := range toRename {
		//fmt.Sprintf("%s/%s", dir, file.Name()),
		//等价于 filepath.Join用slash拼接
		origPath := filepath.Join(dir, origFilename)
		newFilename, err := match(origFilename, count)
		if err != nil {
			panic(err)
		}
		newPath := filepath.Join(dir, newFilename)
		fmt.Printf("mv %s => %s\n", origPath, newPath)
		err = os.Rename(origPath, newPath)
		if err != nil {
			panic(err)
		}
	}
}

//match returns the new file name, or an error if the file name
//didn't match our pattern.
func match(filename string, total int) (string, error) {
	//transfer name like " brithday_001 " to " Brithday - 1 of 4.txt "
	//1. 获取文件名，与后缀分开
	pieces := strings.Split(filename, ".")
	ext := pieces[len(pieces)-1]
	realName := strings.Join(pieces[0:len(pieces)-1], ".")

	//2. 将文件名 按 下划线分割，提取数字
	pieces = strings.Split(realName, "_")
	realName = strings.Join(pieces[0:len(pieces)-1], "_")
	num, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match out pattern!", filename)
	}
	//strings.Title(s)将s中每个单词首字母变为大写，eg:to be number 1 => To Be Number 1
	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(realName), num, total, ext), nil
}
