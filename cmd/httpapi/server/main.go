package main

import (
	"log"
	"net/http"

	"github.com/ccaneke/url-shortner-poc/cmd/httpapi/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", handlers.Shorten)
	mux.HandleFunc("/", handlers.RedirectToLongURL)

	log.Println("Starting server on :3000")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
