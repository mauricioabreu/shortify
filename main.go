package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// Pieces of regular expressions used to validate URLs
const (
	IP           string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	URLSchema    string = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername  string = `(\S+(:\S*)?@)`
	URLPath      string = `((\/|\?|#)[^\s]*)`
	URLPort      string = `(:(\d{1,5}))`
	URLIP        string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	URLSubdomain string = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	URL          string = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
)

var baseURL string
var rgxURL *regexp.Regexp
var redisPool *redis.Client

func init() {
	redisPool = RedisClient()
	rgxURL = regexp.MustCompile(URL)
	flag.StringVar(&baseURL, "base_url", "http://localhost:8000/", "Base site URL")
	flag.Parse()
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Generate a random slug. The final length of this function is twice size
func slugify(url string, size int) string {
	b, err := generateRandomBytes(size)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func shortifyURL(url string) (string, error) {
	err := errors.New("Invalid URL")
	if validateURL(url) != true {
		return "", err
	}

	redisPool := RedisClient()
	shortenedURL := slugify(url, 5)
	err = redisPool.Set(shortenedURL, url, 0).Err()
	if err != nil {
		panic(err)
	}
	return shortenedURL, nil
}

func getURL(slug string) (string, error) {
	redisPool := RedisClient()
	shortenedURL, err := redisPool.Get(slug).Result()
	return shortenedURL, err
}

func shortifyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	shortenedURL, err := shortifyURL(r.Form.Get("url"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return
	}

	log.Printf("Shortened URL: %s", shortenedURL)
	w.Write([]byte(baseURL + shortenedURL))
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

func validateURL(value string) bool {
	if value == "" || strings.HasPrefix(value, ".") {
		return false
	}
	tempValue := value
	// Validate URLs that do not start with a scheme
	if strings.Contains(value, ":") && !strings.Contains(value, "://") {
		tempValue = "http://" + value
	}
	u, err := url.Parse(tempValue)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return rgxURL.MatchString(value)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{slug}", getURLHandler)
	router.HandleFunc("/", shortifyHandler)
	log.Fatal(http.ListenAndServe(":8000", router))
}
