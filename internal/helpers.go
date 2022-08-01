package internal

import (
	"errors"
	"net/url"

	"github.com/google/uuid"
)

func ShortenURL(u url.URL, domain string) (string, error) {
	u.Host = domain
	u.Path = uuid.New().String()[0:6]
	if len(u.Path) > 10 {
		return "", errors.New("ID must be max. 10 characters long")
	}
	rawURL := u.String()
	return rawURL, nil
}
