package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var defaultHandlerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose your own Adventure</title>
    <style>
        body {
            font-family: "Georgia", serif;
            line-height: 1.8;
            margin: 0;
            padding: 0;
            background-color: #fdf6e3;
            color: #333;
        }
        h1 {
            color: #5a3e2b;
            text-align: center;
            margin-top: 40px;
            font-size: 2.5em;
        }
        p {
            margin: 20px 0;
            text-indent: 2em;
            font-size: 1.2em;
        }
        ul {
            list-style: none;
            padding: 0;
            margin-top: 30px;
        }
        li {
            margin: 10px 0;
        }
        a {
            text-decoration: none;
            color: #8b4513;
            font-weight: bold;
        }
        a:hover {
            text-decoration: underline;
            color: #5a3e2b;
        }
        .container {
            max-width: 700px;
            margin: 0 auto;
            padding: 30px;
            background: #fff;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
            border-radius: 10px;
            border: 1px solid #e2d3c5;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Title}}</h1>
        {{range .Paragraphs}}
            <p>{{.}}</p>
        {{end}}
        <ul>
            {{range .Options}}
                <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
            {{end}}
        </ul>
    </div>
</body>
</html>
`

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Story map[string]Chapter

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type handler struct {
	s Story
	t *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		if err := tpl.Execute(w, chapter); err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found", http.StatusNotFound)

}

func NewHandler(s Story, t *template.Template) http.Handler {
	if t == nil {
		t = tpl
	}
	return handler{s, t}
}
