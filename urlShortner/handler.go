package urlshortner

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

//fallback:后备,备用
func MapHandler(pathToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//http.HandlerFunc实现了http.Handler接口
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHander(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	//1. parse yml somehow
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}
	//2. build pathurlmap
	pathToURLs := make(map[string]string)
	for _, pu := range pathUrls {
		pathToURLs[pu.Path] = pu.URL
	}
	//3. call maphandler
	return MapHandler(pathToURLs, fallback), nil
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
