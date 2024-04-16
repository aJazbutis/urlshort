package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"urlshort"
	// "github.com/gophercises/urlshort"
)

func getFileContent(path *string) []byte {
	content, err := os.ReadFile(*path)
	if err != nil {
		panic(err)
	}
	return content
}

func isJson(path *string) bool {
	s := strings.ToLower(*path)
	return strings.HasSuffix(s, ".json")
}

func main() {
	mux := defaultMux()

	var file = flag.String("file", "", "Usage: -file filename.yaml/json")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	content := []byte(yaml)
	if *file != "" {
		content = getFileContent(file)

	}
	fmt.Println("Starting the server on :8080")
	if isJson(file) {
		jsonHandler, err := urlshort.JSONHandler(content, mapHandler)
		if err != nil {
			panic(err)
		}
		http.ListenAndServe(":8080", jsonHandler)

	} else {
		yamlHandler, err := urlshort.YAMLHandler(content, mapHandler)
		if err != nil {
			panic(err)
		}
		http.ListenAndServe(":8080", yamlHandler)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
