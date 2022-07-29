package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", shorten)

	log.Println("Starting server on :3000")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
