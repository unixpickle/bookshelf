package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var GlobalBooksFetcher *BooksFetcher

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: bookshelf <port>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Also, set the GOOGLE_EMAIL and GOOGLE_PASSWORD env vars")
		os.Exit(1)
	}

	email := os.Getenv("GOOGLE_EMAIL")
	password := os.Getenv("GOOGLE_PASSWORD")
	if email == "" || password == "" {
		fmt.Fprintln(os.Stderr, "Missing GOOGLE_EMAIL or GOOGLE_PASSWORD env vars")
		os.Exit(1)
	}

	GlobalBooksFetcher = NewBooksFetcher(email, password)

	http.HandleFunc("/", HandleRoot)
	http.ListenAndServe(":"+os.Args[1], nil)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		books, err := GlobalBooksFetcher.Books()
		var booksData []byte
		if err != nil {
			booksData, _ = json.Marshal(map[string]string{"error": err.Error()})
		} else {
			booksData, _ = json.Marshal(map[string]interface{}{"books": books})
		}
		contents, err := ioutil.ReadFile("assets/index.html")
		if err != nil {
			panic("Failed to read index.html: " + err.Error())
		}
		pageText := strings.Replace(string(contents), "%BOOKS%", string(booksData), 1)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(pageText))
		return
	}

	filename := strings.Replace(path.Base(r.URL.Path), "..", "", -1)
	http.ServeFile(w, r, filepath.Join("assets", filename))
}
