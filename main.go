package main

import (
	"log"
	"net/http"
)

func run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /{shorten}", HandleResolveLink)
	router.HandleFunc("GET /links", HandleGetLinks)
	router.HandleFunc("POST /links/shorten", HandleLinkShorten)

	server := http.Server{Addr: ":8080", Handler: router}
	log.Println("Start server on port 8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server stopped running.", err)
	}

	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
