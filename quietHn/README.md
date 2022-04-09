# **Quite Hack News**

---
### **1. ```http.HandleFunc``` advantages**
1. Normal Usage
   ```go
   hander := func(w, r) {switch r.URL.path() {} }
   http.ListenAndServe(":8080", http.Handler)
   ```
2. With ServeMux
   ```go
   mux := http.NewServeMux()
   mux.Handle(xx1, http.Handler)
   mux.Handle(xx2, http.Handler)
   http.ListenAndServe(":8080", mux)
   ```
3. With DefaultServeMux (global)
   ```go
   http.HandleFunc(xx1, http.Handler)
   http.HandleFunc(xx2, http.Handler)
   http.ListenAndServe(":8080", nil)
   ```
### **2. Usage of ```http.template``` package**
> If you need more complex print formats, you generally need to seperate out the formatting code for safer modification.  

```{{.Time | printf ".2f" }}``` will take the first part as an arguement for the second part.
1. Define a source of template(string or file)
   ```go
   {{range .Stories}}
      <li><a href="{{.URL}}">{{.Title}}</a> <span class="host">({{.Host}})</span></li>
   {{end}}
   <p class="time">This page was rendered in {{.Time | printf "%.2fs" }}</p>
   <p class="footer">This page is heavily inspired by <a href="https://github.com/ldxcwu">Quiet Hacker News</a>.</p>
   ```
2. Create a ```*template.Template```
   ```go
   tpl := template.Must(template.ParseFiles("./index.tmpl"))
   or
   tpl := template.New("tmp").Parse(`string`)
   ```
3. Render
   ```go
   data := Data{Stories: s, URL: url, Host: host, Time: time, xxx}
   err := tpl.Execute(io.Writer, data) 
   ```
### **3. Eliminate competition**
> When we added cache to the project. We have to think about how to eliminate competition when multiple goroutine check the cache at the same time. 

> We can use the command ```go run -race main.go``` to see if there are competitions.   
```go
var (
	cache           []hn.Item
	cacheExpiration time.Time
	mux             sync.Mutex
)

func getCachedStories(numStories int) ([]hn.Item, error) {
	mux.Lock()
	defer mux.Unlock()
	if time.Since(cacheExpiration) < 0 {
		return cache, nil
	}
   // ......
```
---
![image](images/home.jpg)
