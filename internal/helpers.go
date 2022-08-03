package internal

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ccaneke/url-shortner-poc/cmd/httpapi/request"
)

const (
	IDTooLongErrMessage = "id must be max. 10 characters long"
	BlankURLErrMessage  = "long URL cannot be blank"
)

// ShortenURL shortens a long url
func ShortenURL(u url.URL, domain string, uuidString string) (string, error) {
	if len([]rune(uuidString)) > 10 {
		log.Println(IDTooLongErrMessage)
		return "", errors.New(IDTooLongErrMessage)
	}
	u.Host = domain
	u.Path = uuidString
	rawURL := u.String()
	return rawURL, nil
}

// URLFromBody gets the url sent in the body of a request
func URLFromBody(body io.ReadCloser) (*url.URL, error) {
	b, err := io.ReadAll(body)
	if len(b) == 0 {
		log.Println(BlankURLErrMessage)
		return nil, errors.New(BlankURLErrMessage)
	}
	if err != nil {
		log.Println("URLFromBody: ", err)
		return nil, err
	}

	var request request.Request
	err = json.Unmarshal(b, &request)
	if err != nil {
		log.Println("URLFromBody: ", err)
		return nil, err
	}

	defer body.Close()

	u, err := url.Parse(request.LongURL)
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
