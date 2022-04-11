package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/panic", panicDemo)
	mux.HandleFunc("/panic-after", panicAfterDemo)

	http.ListenAndServe(":8080", mux)
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	panic("Oh no!")
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
