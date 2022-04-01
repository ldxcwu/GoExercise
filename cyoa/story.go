package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	var err error
	//按照预先定义的模版去显示内容，
	tpl, err = template.ParseFiles("story.tmpl")
	if err != nil {
		panic(err)
	}
}

type Story map[string]Chapter

//最后的http.ListenAndServe需要接收一个handler，这个handler用以显示story
//因此这里提供一个生成handler的方法，并与story结合起来
func NewHandler(s Story) http.Handler {
	return handler{s}
}

//可以直接为Story实现handler接口，也可以将Story包裹一层
func (s Story) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	if chapter, ok := s[path]; ok {
		if err := tpl.Execute(w, chapter); err != nil {
			log.Fatal(err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "no such page", http.StatusNotFound)
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

//Decoder是解析流（io.Reader)；UnMarshal解析byte切片
func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}
