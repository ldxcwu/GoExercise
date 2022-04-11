# Panic/Recover Middleware
[Check out this article to learn some basic knowledge about defer, panic and recover.](DeferPanicAndRecover.md)

## Panic & Recover
```go

package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/panic", panicDemo)
	mux.HandleFunc("/panic-after", panicAfterDemo)

	http.ListenAndServe(":8080", recoverMiddleware(mux, true))
}

func recoverMiddleware(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if !dev {
					http.Error(w, "Some went wrong.", http.StatusInternalServerError)
					return
				}
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				fmt.Fprintf(w, "<h1>%v</h1><pre>%s</pre>", err, string(stack))
			}
		}()
		app.ServeHTTP(w, r)
	}
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
```

Then:   
1. put codes file on browser to have a test.(with syntax highlighting)   
2. exact filename from panic stack trace  
3. put above source file on browser.
   
## 1. go get github.com/alecthomas/chroma to test view source code on browser.
```go


func sourceCodeHandler(w http.ResponseWriter, r *http.Request) {
	// if path, ok := r.Form["path"]; ok {
	// os.Open(path[0])
	// }
	path := r.FormValue("path")
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	// io.Copy(w, file)
	var bytes bytes.Buffer
	bytes.ReadFrom(file)
	io.Copy(&bytes, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = quick.Highlight(w, bytes.String(), "go", "html", "monokai")
}
