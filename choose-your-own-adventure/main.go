package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func parseFlags() (string, int, string) {
	port := flag.Int("port", 8080, "Port to listen on")
	fileName := flag.String("fileName", "gopher.json", "file to parse")
	htmlTemplatePath := flag.String("templatePath", "", "path to single html template")
	flag.Parse()
	return *fileName, *port, *htmlTemplatePath
}

func buildHtmlTemplate(path string) (*template.Template, error) {

	htmlTemplate, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}
	return htmlTemplate, nil
}

func main() {

	fileName, port, htmlTemplatePath := parseFlags()
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	story, err := JsonStory(file)
	if err != nil {
		log.Fatal(err)
	}

	var options []HandlerOptions

	if htmlTemplatePath != "" {
		tmpl, err := buildHtmlTemplate(htmlTemplatePath)
		if err != nil {
			log.Fatal(err)
		}
		options = append(options, WithHtmlTemplate(tmpl))
	}

	h := NewHandler(story, options...)

	fmt.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), h))
}
