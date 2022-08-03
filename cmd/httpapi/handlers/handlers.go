package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ccaneke/url-shortner-poc/cmd/httpapi/response"
	"github.com/ccaneke/url-shortner-poc/internal"
	db "github.com/ccaneke/url-shortner-poc/internal/DB"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	internalServerError = "Internal Server Error"
)

func Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL, err := internal.URLFromBody(r.Body)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	longURLCopy := *longURL
	domain := internal.GetDomain(r)

	uuidTruncated := uuid.New().String()[0:8]
	shortURL, err := internal.ShortenURL(longURLCopy, domain, uuidTruncated)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	rdb := db.InitRedisDB(ctx)

	err = rdb.Set(ctx, shortURL, longURL.String(), 0).Err()
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	response := response.Response{OriginalURL: longURL.String(), ShortURL: shortURL}
	b, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	w.Write(b)
}

func RedirectToLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	ctx := context.Background()
	rdb := db.InitRedisDB(ctx)

	domain := internal.GetDomain(r)

	longURL, err := getLongURL(ctx, r, rdb, domain)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, *longURL, http.StatusMovedPermanently)
}

type RedisClientInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
}

// getLongURL gets the long url that a short url maps to
func getLongURL(ctx context.Context, r *http.Request, rdb RedisClientInterface, domain string) (*string, error) {

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
