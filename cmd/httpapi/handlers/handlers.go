package handlers

import (
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

type LoggerInterface interface {
	Print(v ...any)
	Fatal(v ...any)
}

type URLShortenerHandler struct {
	redisClient db.RedisClientInterface
	logger      LoggerInterface
}

type ShortnerRedirecter interface {
	Shorten(w http.ResponseWriter, r *http.Request)
	RedirectToLongURL(w http.ResponseWriter, r *http.Request)
}

// NewURLShortenerHandler creates a new URL shortner handler to be used for incoming requests to shorten a url and redirect to a long url from it
func NewURLShortenerHandler(redisClient db.RedisClientInterface, logger LoggerInterface) ShortnerRedirecter {
	return &URLShortenerHandler{redisClient: redisClient, logger: logger}
}

func (h *URLShortenerHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	domainName := getDomainName(r)

	uuidTruncated := uuid.New().String()[0:8]
	shortURL, err := shortenURL(longURLCopy, domainName, uuidTruncated, h.logger)
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

func (h *URLShortenerHandler) RedirectToLongURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	domain := getDomainName(r)

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
func shortenURL(u url.URL, domainName string, uuidString string, logger LoggerInterface) (string, error) {
	if len([]rune(uuidString)) > 10 {
		logger.Print(IDTooLongErrMessage)
		return "", errors.New(IDTooLongErrMessage)
	}
	u.Host = domainName
	u.Path = uuidString
	rawURL := u.String()
	return rawURL, nil
}

// urlFromBody gets the url sent in the body of a request
func urlFromBody(body io.ReadCloser, logger LoggerInterface) (*url.URL, error) {
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

// getDomainName gets the host of the server without the network address
func getDomainName(r *http.Request) string {
	subStrings := strings.Split(r.Host, ":")
	domain := subStrings[0]

	return domain
}
