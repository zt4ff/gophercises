package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"urlshort"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("paths.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	yamlPath := flag.String("y", "", "This takes a path to a YAML file")
	jsonPath := flag.String("j", "", "This takes a path to a JSON FILE")
	flag.Parse()

	pathsToUrls := make(map[string]string)

	if *yamlPath != "" {
		urlshort.ParseYAML(*yamlPath, &pathsToUrls)
	}

	if *jsonPath != "" {
		urlshort.ParseJSON(*jsonPath, &pathsToUrls)
	}

	for path, url := range pathsToUrls {
		urlshort.CreatePath(db, path, url)
	}

	mux := defaultMux(db)

	mapHandler := urlshort.Handler(db, mux)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mapHandler)
}

func defaultMux(db *bolt.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlshort.DefaultHandler)
	mux.HandleFunc("/add-path", urlshort.AddPathHandler(db))
	return mux
}
