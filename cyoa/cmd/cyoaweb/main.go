package main

import (
	"cyoa"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.String("port", "3000", "the port to start the server of the CYOA application")
	filename := flag.String("file", "gopher.json", "The JSON file with the CYOA story")
	flag.Parse()

	fmt.Printf("Using the story in %s\n", *filename)
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	// tpl := template.Must(template.New("").Parse("Hello World!"))
	// h := cyoa.NewHandler(story, cyoa.WithTemplate(tpl))

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server at: %v\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), h))
}

func PathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/storie/intro"
	}

	return path[1:]
}
