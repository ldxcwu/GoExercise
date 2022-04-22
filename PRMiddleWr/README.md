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
```

## 2. extract filename from stack traces  
the stack traces's msg is the specifile pattern,
and we can extract filename from it easily.
```go
// 1. get stack
stack := debug.Stack()
// 2. make links
// goroutine 19 [running]:
// runtime/debug.Stack()
// 	/usr/local/go/src/runtime/debug/stack.go:24 +0x88
// main.recoverMiddleware.func1.1(0x1, {0x100a8bb48, 0x140001ae460})
// 	/Users/lxw/Documents/FTC/GoExercise/PRMiddleWr/main.go:58 +0x74
// panic({0x100a316a0, 0x100a846d8})
func makeLinks(stack string) string {
	lines := strings.Split(stack, "\n")
	for l, line := range lines {
		if len(line) == 0 || line[0] != '\t' {
			continue
		}
		file := ""
		for i, ch := range line {
			if ch == ':' {
				file = line[1:i]
				break
			}
		}
		lines[l] = "\t<a href=\"/debug?path=" + file + "\">" + file + "</a>" + line[len(file)+1:]
	}
	return strings.Join(lines, "\n")
}
//3. encode url path
v := url.Values{}	
v.Set("path", file)
lines[l] = "\t<a href=\"/debug?" + v.Encode() + "\">" + file + "</a>" + line[len(file)+1:]
```
