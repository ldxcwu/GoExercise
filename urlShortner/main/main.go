package main

import (
	"fmt"
	"net/http"
	"urlshortner"
)

func main() {

	pathToURLs := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortner.MapHandler(pathToURLs, http.HandlerFunc(hello))

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshortner.YAMLHander([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	//启动服务器
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello world!")
}
