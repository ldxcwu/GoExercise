package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"link"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*
	1. GET the webpage
	2. parse all the links on the page
	3. build proper urls with our links
	4. filter out any links w/ a diff domain
	5. Find all pages(BFS)
	6. print out XML
*/
func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 3, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	toXml := urlSet{
		Xmlns: xmlns,
	}
	for _, p := range pages {
		toXml.Urls = append(toXml.Urls, loc{p})
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	fmt.Print(xml.Header)
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	Urls []loc `xml:"url"`
	//xml namespace
	Xmlns string `xml:"xmlns,attr"`
}

func bfs(urlStr string, maxDepth int) []string {
	//a set to avoid cyclical避免循环的集合
	seen := make(map[string]struct{})
	//当前层的url
	var q map[string]struct{}
	//下一层的url
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			//已经遍历过
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				if _, ok := seen[link]; !ok {
					nq[link] = struct{}{}
				}
			}
		}
	}
	ret := make([]string, 0)
	for url := range seen {
		ret = append(ret, url)
	}
	return ret

}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	//只用域名，例如url：https://gophercises.com/cyoa/intro
	//则 baseUrl: https://gophercises.com
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	//拼接links中的href
	var hrefs []string
	for _, l := range links {
		switch {
		//link.Parse()解析出来的 a 标签，中的href不一定全是合法的URL，因为可以缩略或mao dian
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	return hrefs
}

//deprecated，这么做过滤规则就写死了，不利于更改，
//事实上应该提供一个方法，根据不同的方法进行不同的过滤
func _filter(links []string, base string) []string {
	var ret []string
	for _, l := range links {
		if strings.HasPrefix(l, base) {
			ret = append(ret, l)
		}
	}
	return ret
}

func filter(links []string, keepFunc func(string) bool) []string {
	var ret []string
	for _, l := range links {
		if keepFunc(l) {
			ret = append(ret, l)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
