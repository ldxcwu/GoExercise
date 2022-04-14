package main

import (
	"cyoa"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 8080, "the port to start the web app")
	filename := flag.String("file", "../../gopher.json", "the JSON file with the CYOA story")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	//%+v先输出字段名字，再输出字段的值
	// fmt.Printf("%+v\n", story)

	fmt.Printf("Starting the server on port: %d\n", *port)
	// http.ListenAndServe(fmt.Sprintf(":%d", *port), story)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), cyoa.NewHandler(story, nil))

}
