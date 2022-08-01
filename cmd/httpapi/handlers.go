package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ccaneke/url-shortner-poc/internal"
	"github.com/go-redis/redis/v8"
)

const (
	separator string = ":"
	fileName  string = "mapping.json"

	internalServerError = "Internal Server Error"
)

func shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u, err := getURL(r.Body)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	copy := *u
	subStrings := strings.Split(r.Host, separator)
	domain := subStrings[0]
	shortenedUrl, err := internal.ShortenURL(copy, domain)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0})

	ctx := context.Background()

	err = rdb.Set(ctx, shortenedUrl, u.String(), 0).Err()
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(shortenedUrl))
}

func getURL(body io.ReadCloser) (*url.URL, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		log.Println("getURL: ", err)
		return nil, err
	}
	defer body.Close()

	u, err := url.Parse(string(b))
	if err != nil {
		log.Println("getURL", err)
		return nil, err
	}

	return u, nil
}
