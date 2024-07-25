package main

import (
	"flag"
	"fmt"
	"htmllp/htmllp"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func parseFlags() (string, string, string) {

	filePath := flag.String("filePath", "files/ex1.html", "file path")
	url := flag.String("url", "", "This flag will ignore the filePath")
	customUrlContains := flag.String("contains", "", "Custom URL contains this flag")

	flag.Parse()

	return *filePath, *url, *customUrlContains
}

func loadReader(filePath, url string) (io.Reader, func()) {

	if url != "" {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		return resp.Body, nil
	}

	fileReader, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %s. %v", filePath, err)
	}
	deferFunc := func() {
		err := fileReader.Close()
		if err != nil {
			log.Fatalf("Error closing file: %s. %v", filePath, err)
		}
	}

	return fileReader, deferFunc
}

func main() {

	filePath, url, customUrlContains := parseFlags()

	reader, deferFunc := loadReader(filePath, url)

	if deferFunc != nil {
		defer deferFunc()
	}

	parser, err := htmllp.NewHtmlParser(reader, func(s string) bool {
		return strings.Contains(s, customUrlContains)
	})

	if err != nil {
		log.Fatalf("Error parsing file to HTML: %s. %v", filePath, err)
	}

	links, err := parser.ReadANodes()
	if err != nil {
		log.Fatalf("Error Reading Nodes: %s. %v", filePath, err)
	}

	fmt.Printf("%+v", links)
}
