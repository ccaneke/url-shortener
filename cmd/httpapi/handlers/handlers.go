package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ccaneke/url-shortner-poc/cmd/httpapi/request"
	"github.com/ccaneke/url-shortner-poc/cmd/httpapi/response"
	db "github.com/ccaneke/url-shortner-poc/internal/DB"
	"github.com/google/uuid"
)

const (
	internalServerError = "Internal Server Error"
	IDTooLongErrMessage = "id must be max. 10 characters long"
	BlankURLErrMessage  = "long URL cannot be blank"
)

type loggerInterface interface {
	Print(v ...any)
	Fatal(v ...any)
}

type handler struct {
	redisClient db.RedisClientInterface
	logger      loggerInterface
}

func NewHandler(redisClient db.RedisClientInterface, logger loggerInterface) handler {
	return handler{redisClient: redisClient, logger: logger}
}

func (h *handler) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL, err := urlFromBody(r.Body, h.logger)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	longURLCopy := *longURL
	domain := getDomain(r)

	uuidTruncated := uuid.New().String()[0:8]
	shortURL, err := shortenURL(longURLCopy, domain, uuidTruncated, h.logger)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	err = h.redisClient.Set(ctx, shortURL, longURL.String(), 0).Err()
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	response := response.Response{OriginalURL: longURL.String(), ShortURL: shortURL}
	b, err := json.Marshal(response)
	if err != nil {
		h.logger.Print(err)
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	w.Write(b)
}

func (h *handler) RedirectToLongURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	domain := getDomain(r)

	r.URL.Host = domain
	r.URL.Scheme = "https"
	key := r.URL.String()
	longURL, err := db.GetValue(ctx, key, h.redisClient, h.logger)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, *longURL, http.StatusMovedPermanently)
}

// shortenURL shortens a long url
func shortenURL(u url.URL, domain string, uuidString string, logger loggerInterface) (string, error) {
	if len([]rune(uuidString)) > 10 {
		logger.Print(IDTooLongErrMessage)
		return "", errors.New(IDTooLongErrMessage)
	}
	u.Host = domain
	u.Path = uuidString
	rawURL := u.String()
	return rawURL, nil
}

// urlFromBody gets the url sent in the body of a request
func urlFromBody(body io.ReadCloser, logger loggerInterface) (*url.URL, error) {
	b, err := io.ReadAll(body)
	if len(b) == 0 {
		logger.Print(BlankURLErrMessage)
		return nil, errors.New(BlankURLErrMessage)
	}
	if err != nil {
		logger.Print("URLFromBody: ", err)
		return nil, err
	}

	var request request.Request
	err = json.Unmarshal(b, &request)
	if err != nil {
		logger.Print("URLFromBody: ", err)
		return nil, err
	}

	defer body.Close()

	u, err := url.Parse(request.LongURL)
	if err != nil {
		logger.Print("getURL", err)
		return nil, err
	}

	return u, nil
}

// getDomain gets the host of the server without the network address
func getDomain(r *http.Request) string {
	subStrings := strings.Split(r.Host, ":")
	domain := subStrings[0]

	return domain
}
