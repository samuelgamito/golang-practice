package htmllp

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type HtmlLinkParser struct {
	RootNode *html.Node
	Filter   FilterOption
}

type Link struct {
	Url  string
	Text string
}

type FilterOption func(string) bool

func defaultFilter(s string) bool {
	return true
}

func NewHtmlParser(ioReader io.Reader, filter FilterOption) (*HtmlLinkParser, error) {
	doc, err := html.Parse(ioReader)
	if err != nil {
		return nil, err
	}

	parser := &HtmlLinkParser{RootNode: doc, Filter: defaultFilter}

	if filter != nil {
		parser.Filter = filter
	}

	return parser, nil
}

func (h HtmlLinkParser) ReadANodes() ([]Link, error) {

	var f func(*html.Node)

	var links []Link
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			if n.Attr != nil {
				for _, a := range n.Attr {
					if a.Key == "href" && h.Filter(a.Val) {
						links = append(links, Link{a.Val, strings.TrimSpace(n.FirstChild.Data)})
						break
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(h.RootNode)

	return links, nil
}