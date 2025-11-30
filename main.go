package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func run() error {
	ctx := context.Background()
	dsn := os.Getenv("PG_ADDR")
	fmt.Printf("ENV: %s \n", dsn)

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Fatal("Could not connect to database.", err)
	}

	defer pool.Close()

	repo := LinksRepository{Pool: pool}

	router := http.NewServeMux()
	router.HandleFunc("GET /{shorten}", HandleResolveLink(&repo))
	router.HandleFunc("GET /links", HandleGetLinks(&repo))
	router.HandleFunc("POST /links/shorten", HandleLinkShorten(&repo))

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
