package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
)

type AddPathRequest struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func AddPathHandler(db *bolt.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req AddPathRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Path == "" || req.URL == "" {
			http.Error(w, "Path and URL must not be empty", http.StatusBadRequest)
			return
		}

		err := CreatePath(db, req.Path, req.URL)
		if err != nil {
			http.Error(w, "Failed to save path", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Path added successfully"))
	}
}

func Handler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		var url string

		GetPath(db, path, &url)

		if url != "" {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You may have inputed an invalid path")
}
