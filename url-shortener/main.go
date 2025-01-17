package main

import (
	urlshort "exercise/urlshort/handler"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	mux := defaultMux()
	file := flag.String("f", "", "The file with the paths and urls")
	flag.Parse()

	var fileContent []byte
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	fileContent = []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	if *file != "" {
		file, err := os.ReadFile(*file)
		fileContent = file
		checkError(err)
	}

	var handler http.HandlerFunc

	if strings.Split(*file, ".")[1] == "json" {
		jsonHandler, err := urlshort.JSONHandler(fileContent, mapHandler)
		handler = jsonHandler
		checkError(err)

	} else if strings.Split(*file, ".")[1] == "yml" {
		yamlHandler, err := urlshort.YAMLHandler(fileContent, mapHandler)
		handler = yamlHandler
		checkError(err)
	} else {
		fmt.Println("Invalid file format")
		os.Exit(1)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
