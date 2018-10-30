package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func shortifyURL(url string) string {
	redis := RedisClient()
	shortenedURL := "abc012"
	err := redis.Set(shortenedURL, url, 0).Err()
	if err != nil {
		panic(err)
	}
	return shortenedURL
}

func getURL(slug string) (string, error) {
	redis := RedisClient()
	shortenedURL, err := redis.Get(slug).Result()
	return shortenedURL, err
}

func shortifyHandler(w http.ResponseWriter, r *http.Request) {
	shortenedURL := shortifyURL("foo")
	log.Printf("Shortened URL: %s", shortenedURL)
	fmt.Fprintf(w, shortenedURL)
}

func getURLHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shortenedURL, err := getURL(params["slug"])
	if shortenedURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	} else {
		w.Header().Add("Location", shortenedURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{slug}", getURLHandler)
	router.HandleFunc("/", shortifyHandler)
	log.Fatal(http.ListenAndServe(":8000", router))
}
