package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ccaneke/url-shortner-poc/internal"
)

const (
	colon string = ":"
)

func shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println("url: ", r.Host)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	u, err := url.Parse(string(body))
	copy := *u
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s := strings.Split(r.Host, colon)
	domain := s[0]
	shortenedUrl, err := internal.ShortenURL(copy, domain)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("url from request body: %v, shortened url: %v", u, shortenedUrl)

	w.WriteHeader(http.StatusOK)
}
