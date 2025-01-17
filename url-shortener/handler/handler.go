package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
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
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func buildMap(pathsToUrls []pathURL) map[string]string {
	paths := make(map[string]string)
	for _, pu := range pathsToUrls {
		paths[pu.Path] = pu.URL
	}
	return paths
}

func parseYAML(yml []byte) ([]pathURL, error) {
	var paths []pathURL
	err := yaml.Unmarshal(yml, &paths)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func parseJSON(jsonData []byte) ([]pathURL, error) {
	var paths []pathURL
	if err := json.Unmarshal(jsonData, &paths); err != nil {
		return nil, err
	}

	return paths, nil
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseJSON(jsonData)
	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(paths), fallback), nil
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	paths, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(paths), fallback), nil
}
