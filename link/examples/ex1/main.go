package main

import (
	"flag"
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
	_, err = link.Parse(f)
	if err != nil {
		panic(err)
	}
}
