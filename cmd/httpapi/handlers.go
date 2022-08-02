package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ccaneke/url-shortner-poc/internal"
	"github.com/go-redis/redis/v8"
)

const (
	internalServerError = "Internal Server Error"
)

func shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL, err := internal.GetURL(r.Body)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	longURLCopy := *longURL
	domain := internal.GetDomain(r)

	shortURL, err := internal.ShortenURL(longURLCopy, domain)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	rdb, ctx := initRedisDB()

	err = rdb.Set(ctx, shortURL, longURL.String(), 0).Err()
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(shortURL))
}

// initRedisDB connects to a redis server
func initRedisDB() (*redis.Client, context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0})

	ctx := context.Background()

	return rdb, ctx
}

func showLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	rdb, ctx := initRedisDB()

	domain := internal.GetDomain(r)

	longURL, err := getLongURL(ctx, r, rdb, domain)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, *longURL, http.StatusFound)
}

// getLongURL gets the long url that a short url maps to
func getLongURL(ctx context.Context, r *http.Request, rdb *redis.Client, domain string) (*string, error) {
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
