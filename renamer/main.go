package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

type file struct {
	name string
	path string
}

func main() {
	dir := "./sample"

	var toRename []file
	//filepath.Walk(dir, func)会遍历目录下所有文件（不论是文件对象还是目录对象）
	//并对所有对象调用func方法
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		//如果文件是目录则不做任何处理（并非跳过其内部 = filepath.SkipDir())
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, file{
				name: info.Name(),
				path: path,
			})
		}
		return nil
	})
	for _, f := range toRename {
		fmt.Printf("%q\n", f)
	}
	for _, orig := range toRename {
		var newFile file
		var err error
		newFile.name, err = match(orig.name)
		if err != nil {
			fmt.Println("Error matching:", orig.path, err.Error())
		}
		newFile.path = filepath.Join(dir, newFile.name)
		fmt.Printf("mv %s => %s\n", orig.path, newFile.path)
		// err = os.Rename(orig.path, newFile.path)
		// if err != nil {
		// 	fmt.Println("Error renaming:", orig.path, err.Error())
		// }
	}

	// files, err := ioutil.ReadDir(dir)
	// if err != nil {
	// 	fmt.Println("ioutil.ReadDir error")
	// 	os.Exit(1)
	// }
	// var toRename []string
	// count := 0
	// for _, file := range files {
	// 	if !file.IsDir() {
	// 		//如果是文件，就获取文件名并进行match
	// 		_, err := match(file.Name(), 4)
	// 		if err == nil {
	// 			count++
	// 			toRename = append(toRename, file.Name())
	// 		}
	// 	}
	// }
	// for _, origFilename := range toRename {
	// 	//fmt.Sprintf("%s/%s", dir, file.Name()),
	// 	//等价于 filepath.Join用slash拼接
	// 	origPath := filepath.Join(dir, origFilename)
	// 	newFilename, err := match(origFilename, count)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	newPath := filepath.Join(dir, newFilename)
	// 	fmt.Printf("mv %s => %s\n", origPath, newPath)
	// 	err = os.Rename(origPath, newPath)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}

//match returns the new file name, or an error if the file name
//didn't match our pattern.
func match(filename string) (string, error) {
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
	return fmt.Sprintf("%s - %d.%s", strings.Title(realName), num, ext), nil
}
