package main

import (
	"flag"
	"fmt"
	"htmllp"
	"log"
	"net/http"
	"strings"
)

func parseFlags() string {

	filePath := flag.String("filePath", "files/ex1.html", "file path")
	flag.Parse()

	return *filePath
}

func main() {

	filePath := parseFlags()

	//fileReader, err := os.Open(filePath)
	//if err != nil {
	//	log.Fatalf("Error opening file: %s. %v", filePath, err)
	//}
	//defer func() {
	//	err := fileReader.Close()
	//	if err != nil {
	//		log.Fatalf("Error closing file: %s. %v", filePath, err)
	//	}
	//}()

	resp, err := http.Get("https://www.calhoun.io/")

	if err != nil {
		log.Fatal(err)
	}

	parser, err := htmllp.NewHtmlParser(resp.Body, func(s string) bool {
		return strings.Contains(s, "calhoun.io")
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
