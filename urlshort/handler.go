package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type pathUrl struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func contentToMap(content []byte, isYaml bool) (map[string]string, error) {
	data := []pathUrl{}
	if isYaml {
		if err := yaml.Unmarshal(content, &data); err != nil {
			return nil, err
		}
	} else {
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, err
		}
	}
	pathMap := make(map[string]string)
	for _, entry := range data {
		pathMap[entry.Path] = entry.Url
	}
	return pathMap, nil
}

func YAMLHandler(content []byte, fallback http.Handler) (http.HandlerFunc, error) {
	if pathMap, err := contentToMap(content, true); err == nil {
		return MapHandler(pathMap, fallback), nil
	} else {
		return nil, err
	}
}

func JSONHandler(content []byte, fallback http.Handler) (http.HandlerFunc, error) {
	if pathMap, err := contentToMap(content, false); err == nil {
		return MapHandler(pathMap, fallback), nil
	} else {
		return nil, err
	}
}
