package internal

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

// ShortenURL shortens a long url
func ShortenURL(u url.URL, domain string) (string, error) {
	u.Host = domain
	u.Path = uuid.New().String()[0:6]
	if len(u.Path) > 10 {
		return "", errors.New("ID must be max. 10 characters long")
	}
	rawURL := u.String()
	return rawURL, nil
}

// GetURL gets the url sent in the body of a request
func GetURL(body io.ReadCloser) (*url.URL, error) {
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

// GetDomain gets the host of the server without the network address
func GetDomain(r *http.Request) string {
	subStrings := strings.Split(r.Host, ":")
	domain := subStrings[0]

	return domain
}
