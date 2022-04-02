package main

import (
	"flag"
	"fmt"
	"link"
	"os"
)

func main() {
	filename := flag.String("file", "ex1.html", "HTML document you want to parse")
	flag.Parse()
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	links, err := link.Parse(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(links)
}
