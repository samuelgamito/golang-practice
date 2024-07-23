package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("story").Parse(defaultHandlerTemplate))
}

var defaultHandlerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Choose Your Own Adventure</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
            color: #333;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            text-align: center;
        }
        .container {
            background-color: #fff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            width: 100%;
            margin: 20px;
        }
        h1 {
            color: #2c3e50;
        }
        p {
            line-height: 1.6;
        }
        ul {
            list-style: none;
            padding: 0;
        }
        li {
            margin: 10px 0;
        }
        a {
            text-decoration: none;
            color: #3498db;
            font-weight: bold;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
        <p>{{.}}</p>
        {{end}}
		<hr />
        <ul>
            {{range .Options}}
            <li>
                <a href="/{{.NextChapter}}">{{.Text}}</a>
            </li>
            {{end}}
        </ul>
    </div>
</body>
</html>
`

type Story map[string]ArcNode

type ArcOptions struct {
	Text        string `json:"text"`
	NextChapter string `json:"arc"`
}

type ArcNode struct {
	Title      string       `json:"title"`
	Paragraphs []string     `json:"story"`
	Options    []ArcOptions `json:"options"`
}

func JsonStory(reader io.Reader) (Story, error) {

	decoder := json.NewDecoder(reader)
	var story Story
	err := decoder.Decode(&story)
	if err != nil {
		return nil, err
	}
	return story, nil
}

type HandlerOptions func(h *handler)

func WithHtmlTemplate(tmpl *template.Template) HandlerOptions {
	return func(h *handler) {
		h.t = tmpl
	}
}

func NewHandler(s Story, opts ...HandlerOptions) http.Handler {
	h := handler{s, tmpl}

	for _, fn := range opts {
		fn(&h)
	}

	return &h
}

type handler struct {
	story Story
	t     *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]

	if chapter, ok := h.story[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	http.NotFound(w, r)

}
