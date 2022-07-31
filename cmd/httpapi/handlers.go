package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/ccaneke/url-shortner-poc/internal"
)

const (
	separator string = ":"
	fileName  string = "mapping.json"

	internalServerError = "Internal Server Error"
)

var mapping map[string]string = make(map[string]string)

func shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	u, err := getURL(w, r.Body)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	copy := *u
	subStrings := strings.Split(r.Host, separator)
	domain := subStrings[0]
	shortenedUrl, err := internal.ShortenURL(copy, domain)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mapping[shortenedUrl] = u.String()

	projectPath, err := os.Getwd()
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	absFilePath := projectPath + "/" + fileName
	if _, err := os.Stat(absFilePath); errors.Is(err, os.ErrNotExist) {
		err = save(absFilePath, mapping)
		if err != nil {
			http.Error(w, internalServerError, http.StatusInternalServerError)
			return
		}
		return
	}

	m, err := sync(absFilePath, mapping)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	err = save(absFilePath, m)
	if err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(shortenedUrl))
}

func save(file string, changes map[string]string) error {
	b, err := json.Marshal(changes)
	if err != nil {
		log.Println("save:", err)
		return err
	}

	err = os.WriteFile(file, b, 0755)
	if err != nil {
		log.Println("save:", err)
		return err
	}

	return nil
}

func sync(file string, changes map[string]string) (map[string]string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		log.Println("sync:", err)
		return nil, err
	}

	var m map[string]string

	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Println("sync: ", err)
		return nil, err
	}

	maps.Copy(m, changes)

	return m, nil
}

func getURL(w http.ResponseWriter, body io.ReadCloser) (*url.URL, error) {
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
