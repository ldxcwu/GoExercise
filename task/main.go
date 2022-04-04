package main

import (
	"fmt"
	"os"
	"path/filepath"
	"task/cmd"
	"task/db"

	"github.com/mitchellh/go-homedir"
)

func main() {
	//homedir.Dir()返回用户目录，也可以使用os.UserHomeDir()实现，
	//但是那依赖于Darwin系统，无法做到交叉编译，例如在Windows系统上面，
	//所以这里引用这个包，内部做了一下判断
	home, _ := homedir.Dir()
	// fmt.Printf("homeDir: %v\n", home)
	dbPath := filepath.Join(home, "tasks.db")
	must(db.Init(dbPath))
	must(cmd.RootCmd.Execute())

}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
