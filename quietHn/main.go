package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"quietHn/hn"
	"strings"
	"time"
)

func main() {
	port := flag.Int("port", 8080, "The port to start the web server on")
	numStories := flag.Int("num", 30, "The number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.tmpl"))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler(*numStories, tpl)))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var c hn.Client
		ids, err := c.GetTopStories()
		if err != nil {
			http.Error(w, "Failed to load top stories.", http.StatusInternalServerError)
			return
		}
		var td templateData
		for _, id := range ids {
			item, err := c.GetItem(id)
			if err != nil {
				continue
			}
			if item.Type == "story" && item.URL != "" {
				item.Host = url2Host(item.URL)
				td.Stories = append(td.Stories, item)
			}
			if len(td.Stories) > numStories {
				break
			}
		}
		td.Time = time.Since(start).Seconds()
		err = tpl.Execute(w, td)
		if err != nil {
			http.Error(w, "Failed to process the template.", http.StatusInternalServerError)
		}
	})
}

func url2Host(URL string) string {
	u, err := url.Parse(URL)
	if err == nil {
		return strings.TrimPrefix(u.Hostname(), "www.")
	}
	return ""
}

type templateData struct {
	Stories []hn.Item
	Time    float64
}
