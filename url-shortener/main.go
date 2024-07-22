package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func parseFlags() (int, string, string) {
	port := flag.Int("port", 8080, "Port to listen on --port")
	yamlFilePath := flag.String("yaml", "", "Path to YAML config file")
	jsonFilePath := flag.String("json", "", "Path to JSON config file")

	flag.Parse()
	return *port, *yamlFilePath, *jsonFilePath
}

func main() {

	port, yamlFilePath, jsonFilePath := parseFlags()
	mux := defaultMux()

	handlerFuncMux := buildFallbackMap(mux)
	handlerFuncMux = buildYamlHandler(yamlFilePath, handlerFuncMux)

	if jsonFilePath != "" {
		handlerFuncMux = buildJsonHandler(jsonFilePath, handlerFuncMux)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handlerFuncMux)
	if err != nil {
		log.Fatal(err)
	}
}

func buildJsonHandler(path string, mux http.HandlerFunc) http.HandlerFunc {

	fileData, err := os.ReadFile(path)
	var jsonBytes []byte

	if err != nil {
		log.Printf("Error reading file(%s): %v\n", path, err)
		return mux
	} else {
		jsonBytes = fileData
	}

	handler, err := JSONHandler(jsonBytes, mux)

	if err != nil {
		log.Printf("Not able do decode json. %v\n", err)
		return mux
	}

	log.Printf("Appending %s json file to the map\n", path)
	return handler
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	return mux
}

func buildYamlHandler(path string, handlerFuncMux http.HandlerFunc) http.HandlerFunc {
	var yamlBytes []byte

	if path != "" {
		fileData, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file(%s): %v\n", fileData, err)
		} else {
			yamlBytes = fileData
		}

	}

	if yamlBytes == nil {
		log.Println("Using the default YAML Mapping")
		yamlBytes = []byte(
			`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	}

	handler, err := YAMLHandler(yamlBytes, handlerFuncMux)

	if err != nil {
		log.Fatal(err)
	}

	return handler
}
func buildFallbackMap(mux *http.ServeMux) http.HandlerFunc {
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	return MapHandler(pathsToUrls, mux)
}

func ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Print(err)
		return
	}
}
