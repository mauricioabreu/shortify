// These tests ensure that our HTTP handlers are working as expected
package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestShortifyURL(t *testing.T) {
	data := url.Values{}
	data.Set("url", "http://www.maugzoide.com/")
	req, err := http.NewRequest("POST", "localhost:8000/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}
	recorder := httptest.NewRecorder()
	shortifyURLHandler(recorder, req)

	result := recorder.Result()
	defer result.Body.Close()

	b, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP %v, got %v", http.StatusOK, result.StatusCode)
	}

	body := string(bytes.TrimSpace(b))
	slug := strings.Split(body, "/")[3]

	// Ensure we have the slug generated and it is 10 chars length
	if len(slug) != 10 {
		t.Errorf("Slug must be 10 length. Got %v", len(slug))
	}
}

func TestShortifyInvalidURL(t *testing.T) {
	data := url.Values{}
	data.Set("url", "http:invalid-url.com/")
	req, err := http.NewRequest("POST", "localhost:8000/", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}
	recorder := httptest.NewRecorder()
	shortifyURLHandler(recorder, req)

	result := recorder.Result()
	defer result.Body.Close()

	b, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP %v, got %v", http.StatusBadRequest, result.StatusCode)
	}

	body := string(bytes.TrimSpace(b))
	if body != "Invalid URL" {
		t.Errorf("Body content should be 'Invalid URL'. Got %v", body)
	}
}

func TestGetURL(t *testing.T) {
	redisPool := RedisClient()
	err := redisPool.Set("foobar1234", "https://www.maugzoide.com/", 0).Err()
	if err != nil {
		log.Fatalf("Could not set a new key on store: %v", err)
	}

	req, err := http.NewRequest("GET", "localhost:8000/foobar1234", nil)
	// Set a value for the slug path placeholder
	// See https://godoc.org/github.com/gorilla/mux#SetURLVars for more information
	req = mux.SetURLVars(req, map[string]string{"slug": "foobar1234"})
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}
	recorder := httptest.NewRecorder()
	getURLHandler(recorder, req)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Expected HTTP %v, got %v", http.StatusTemporaryRedirect, result.StatusCode)
	}
}
