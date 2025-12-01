package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run() error {
	ctx := context.Background()
	dsn := os.Getenv("PG_ADDR")

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Fatal("Could not connect to database.", err)
	}

	defer pool.Close()

	repo := LinksRepository{pool}
	service := LinksService{&repo}

	router := http.NewServeMux()
	router.HandleFunc("GET /{shorten}", HandleResolveLink(&service))
	router.HandleFunc("GET /links", HandleGetLinks(&service))
	router.HandleFunc("POST /links/shorten", HandleLinkShorten(&service))

	server := http.Server{Addr: ":8080", Handler: router}
	log.Println("Start server on port 8080")
	err = server.ListenAndServe()

	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
