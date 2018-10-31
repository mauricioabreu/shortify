package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var baseURL string

func init() {
	flag.StringVar(&baseURL, "base_url", "http://localhost:8000/", "Base site URL")
	flag.Parse()
}

func slugify(url string) string {
	hash := md5.New()
	return hex.EncodeToString(hash.Sum([]byte(url)))[0:7]
}

func shortifyURL(url string) string {
	redis := RedisClient()
	shortenedURL := slugify(url)
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
	r.ParseForm()
	shortenedURL := shortifyURL(r.Form.Get("url"))
	log.Printf("Shortened URL: %s", shortenedURL)
	fmt.Fprintf(w, shortenedURL)
}

func getURLHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	url, err := getURL(params["slug"])
	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		panic(err)
	} else {
		w.Header().Add("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{slug}", getURLHandler)
	router.HandleFunc("/", shortifyHandler)
	log.Fatal(http.ListenAndServe(":8000", router))
}
