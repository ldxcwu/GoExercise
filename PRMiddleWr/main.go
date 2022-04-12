package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/panic", panicDemo)
	mux.HandleFunc("/panic-after", panicAfterDemo)
	mux.HandleFunc("/debug", sourceCodeHandler)

	http.ListenAndServe(":8080", recoverMiddleware(mux, true))
	// http.ListenAndServe(":8080", mux)
}

/* goroutine 19 [running]:
runtime/debug.Stack()
	/usr/local/go/src/runtime/debug/stack.go:24 +0x88
main.recoverMiddleware.func1.1(0x1, {0x100a8bb48, 0x140001ae460})
	/Users/lxw/Documents/FTC/GoExercise/PRMiddleWr/main.go:58 +0x74
panic({0x100a316a0, 0x100a846d8})
	/usr/local/go/src/runtime/panic.go:1038 +0x21c
main.panicDemo({0x100a8bb48, 0x140001ae460}, 0x140001a2300)
	/Users/lxw/Documents/FTC/GoExercise/PRMiddleWr/main.go:106 +0x38
net/http.HandlerFunc.ServeHTTP(0x100a838f8, {0x100a8bb48, 0x140001ae460}, 0x140001a2300)
	/usr/local/go/src/net/http/server.go:2046 +0x40
net/http.(*ServeMux).ServeHTTP(0x140001fdf80, {0x100a8bb48, 0x140001ae460}, 0x140001a2300)
	/usr/local/go/src/net/http/server.go:2424 +0x18c
main.recoverMiddleware.func1({0x100a8bb48, 0x140001ae460}, 0x140001a2300)
	/Users/lxw/Documents/FTC/GoExercise/PRMiddleWr/main.go:66 +0x7c
net/http.HandlerFunc.ServeHTTP(0x140002728e0, {0x100a8bb48, 0x140001ae460}, 0x140001a2300)
	/usr/local/go/src/net/http/server.go:2046 +0x40
net/http.serverHandler.ServeHTTP({0x140001ae0e0}, {0x100a8bb48, 0x140001ae460}, 0x140001a2300)
	/usr/local/go/src/net/http/server.go:2878 +0x444
net/http.(*conn).serve(0x14000136dc0, {0x100a8ce80, 0x14000231110})
	/usr/local/go/src/net/http/server.go:1929 +0xb6c
created by net/http.(*Server).Serve
	/usr/local/go/src/net/http/server.go:3033 +0x4b8 */

// Parse stack string with links.
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
		v := url.Values{}
		v.Set("path", file)
		lines[l] = "\t<a href=\"/debug?" + v.Encode() + "\">" + file + "</a>" + line[len(file)+1:]
		// lines[l] = "\t<a href=\"/debug?path=" + file + "\">" + file + "</a>" + line[len(file)+1:]
	}
	return strings.Join(lines, "\n")
}

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
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>%v</h1><pre>%s</pre>", err, makeLinks(string(stack)))
			}
		}()
		// rw := &responseWriter{ResponseWriter: w}
		// app.ServeHTTP(rw, r)
		app.ServeHTTP(w, r)
		// rw.flush()
	}
}

// type ResponseWriter interface {
// 	Header() Header
// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

type responseWriter struct {
	http.ResponseWriter
	status int
	writes [][]byte
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
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
