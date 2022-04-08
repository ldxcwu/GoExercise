package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"quietHn/hn"
	"sort"
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

		td, err := getTopStories(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		td.Time = time.Since(start).Seconds()
		err = tpl.Execute(w, td)
		if err != nil {
			http.Error(w, "Failed to process the template.", http.StatusInternalServerError)
			return
		}
	})
}

func getTopStories(numStories int) (templateData, error) {
	var c hn.Client
	var td templateData
	ids, err := c.GetTopStories()
	if err != nil {
		return td, err
	}
	//make sure that we got the correct number of stories.
	var items []hn.Item
	at := 0
	for len(items) < numStories && at < len(ids) {
		need := (numStories - len(items)) * 5 / 4
		//TODO: make sure that at + need < bounds.
		items = append(items, getStories(c, ids[at:at+need])...)
		at += need
	}
	td.Stories = items[0:numStories]
	return td, nil
}

type ret struct {
	idx  int
	item hn.Item
	err  error
}

func getStories(c hn.Client, ids []int) []hn.Item {
	retCh := make(chan ret)
	for i := 0; i < len(ids); i++ {
		go func(i, id int) {
			item, err := c.GetItem(id)
			if err != nil {
				retCh <- ret{idx: i, err: err}
			} else {
				retCh <- ret{idx: i, item: item}
			}
		}(i, ids[i])
	}
	var res []ret
	for i := 0; i < len(ids); i++ {
		res = append(res, <-retCh)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].idx <= res[j].idx
	})
	var items []hn.Item
	for _, r := range res {
		if r.err != nil {
			continue
		}
		item := r.item
		if item.Type == "story" && item.URL != "" {
			item.Host = url2Host(item.URL)
			items = append(items, item)
		}
	}
	close(retCh)
	return items
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
