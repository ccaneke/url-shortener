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
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	rdb, ctx := initRedisDB()

	err = rdb.Set(ctx, shortenedUrl, u.String(), 0).Err()
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(shortenedUrl))
}

func initRedisDB() (*redis.Client, context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0})

	ctx := context.Background()

	return rdb, ctx
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

func showLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	rdb, ctx := initRedisDB()

	subStrings := strings.Split(r.Host, separator)
	domain := subStrings[0]

	longURL, err := getLongURL(ctx, r, rdb, domain)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, *longURL, http.StatusFound)
}

type RedisClient interface {
	Get(context.Context, string) *redis.StringCmd
}

func getLongURL(ctx context.Context, r *http.Request, rdb RedisClient, domain string) (*string, error) {
	r.URL.Host = domain
	r.URL.Scheme = "https"
	key := r.URL.String()

	val, err := rdb.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		log.Println("getLongURL: key does not exist")
		return nil, err
	case err != nil:
		log.Println("getLongURL: Get failed")
		return nil, err
	case val == "":
		log.Println("getLongURL: value is empty")
		return &val, nil
	}

	return &val, nil
}
