package main

import (
	"fmt"
	"log"
	"net/http"
)

func shortifyURL(url string) string {
	return "abc012"
}

func shortifyHandler(w http.ResponseWriter, r *http.Request) {
	shortenedURL := shortifyURL("foo")
	log.Printf("Shortened URL: %s", shortenedURL)
	fmt.Fprintf(w, shortenedURL)
}

func main() {
	http.HandleFunc("/", shortifyHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
