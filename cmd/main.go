package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CipherN9/go-url-shortener/internal/short-links"
	"github.com/jackc/pgx/v5/pgxpool"
)

func importRoutes(router *http.ServeMux, pool *pgxpool.Pool) {
	repo := short_links.LinksRepository{Pool: pool}
	service := short_links.LinksService{Repo: &repo}

	router.HandleFunc("GET /{shorten}", short_links.HandleResolveLink(&service))
	router.HandleFunc("GET /links", short_links.HandleGetLinks(&service))
	router.HandleFunc("POST /links/shorten", short_links.HandleLinkShorten(&service))
}

func run(ctx context.Context) error {
	dsn := os.Getenv("PG_ADDR")

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Fatal("Could not create pool", err)
	}

	poolCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(poolCtx); err != nil {
		{
			return fmt.Errorf("сould not connect to database. %w", err)
		}
	}

	defer pool.Close()

	router := http.NewServeMux()
	importRoutes(router, pool)

	server := http.Server{Addr: ":8080", Handler: router}
	log.Println("Start server on port 8080")

	go func(ctx context.Context) {
		<-ctx.Done()
		log.Println("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server Shutdown error: %v", err)
		}

	}(ctx)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	log.Println("Server shutdown successfully")
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}
