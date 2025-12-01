package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/CipherN9/go-url-shortener/modules/shorten-links"
	"github.com/jackc/pgx/v5/pgxpool"
)

func importRoutes(router *http.ServeMux, pool *pgxpool.Pool) {
	repo := shorten_links.LinksRepository{pool}
	service := shorten_links.LinksService{&repo}

	router.HandleFunc("GET /{shorten}", shorten_links.HandleResolveLink(&service))
	router.HandleFunc("GET /links", shorten_links.HandleGetLinks(&service))
	router.HandleFunc("POST /links/shorten", shorten_links.HandleLinkShorten(&service))
}

func run() error {
	ctx := context.Background()
	dsn := os.Getenv("PG_ADDR")

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Fatal("Could not connect to database.", err)
	}

	defer pool.Close()

	router := http.NewServeMux()
	importRoutes(router, pool)

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
