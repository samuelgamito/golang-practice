package main

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"net/http"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		}

		fallback.ServeHTTP(w, r)
	}
}

func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var parsedJson []map[string]string
	if err := json.Unmarshal(jsonBytes, &parsedJson); err != nil {
		return nil, err
	}

	parsedJsonMap := make(map[string]string)
	for _, pathMap := range parsedJson {
		path := pathMap["path"]
		url := pathMap["url"]

		parsedJsonMap[path] = url
	}

	return MapHandler(parsedJsonMap, fallback), nil
}

func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var parsedYaml interface{}
	parsedYamlMap := make(map[string]string)

	if err := yaml.Unmarshal(yamlBytes, &parsedYaml); err != nil {
		return nil, err
	}

	for _, pathMap := range parsedYaml.([]interface{}) {
		path := pathMap.(map[string]interface{})["path"].(string)
		url := pathMap.(map[string]interface{})["url"].(string)

		parsedYamlMap[path] = url
	}
	return MapHandler(parsedYamlMap, fallback), nil
}
