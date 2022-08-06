package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ccaneke/url-shortner-poc/cmd/httpapi/handlers"
	db "github.com/ccaneke/url-shortner-poc/internal/DB"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)

	redisClient := db.InitRedisDB(logger)
	handler := handlers.NewHandler(redisClient, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", handler.Shorten)
	mux.HandleFunc("/", handler.RedirectToLongURL)

	logger.Print("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	logger.Fatal(err)
}
