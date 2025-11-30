package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ILinksRepository interface {
	Add(context.Context, *Link) error
	Get(context.Context, *Filter) ([]Link, error)
}
type LinksRepository struct {
	Pool *pgxpool.Pool
}

type Link struct {
	Id          uuid.UUID
	InitialLink string
	ShortenLink string
}

type Filter Link

func (r *LinksRepository) Add(ctx context.Context, l *Link) error {
	_, err := r.Pool.Exec(ctx, `INSERT INTO links (id, initial_link, shorten_link) VALUES ($1, $2, $3)`,
		uuid.New(), l.InitialLink, l.ShortenLink)

	if err != nil {
		log.Printf("Insert failed: %v", err)
	}

	return err
}

func (r *LinksRepository) Get(ctx context.Context, filter *Filter) ([]Link, error) {
	baseQuery := `SELECT id, initial_link, shorten_link FROM links`

	var conditions []string
	var args []any

	if filter != nil {
		argPos := 1

		if filter.Id != uuid.Nil {
			conditions = append(conditions, fmt.Sprintf("id = $%d", argPos))
			args = append(args, filter.Id)
			argPos++
		}

		if filter.InitialLink != "" {
			conditions = append(conditions, fmt.Sprintf("initial_link = $%d", argPos))
			args = append(args, filter.InitialLink)
			argPos += 1
		}

		if filter.ShortenLink != "" {
			conditions = append(conditions, fmt.Sprintf("shorten_link = $%d", argPos))
			args = append(args, filter.ShortenLink)
			argPos += 1
		}

		if len(conditions) > 0 {
			baseQuery = baseQuery + " WHERE " + strings.Join(conditions, " AND ")
		}

	}

	rows, err := r.Pool.Query(ctx, baseQuery, args...)

	if err != nil {
		return nil, fmt.Errorf("links query failed: %w", err)
	}

	defer rows.Close()

	var result []Link

	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.Id, &l.InitialLink, &l.ShortenLink); err != nil {
			return nil, fmt.Errorf("links scan failed: %w", err)
		}
		result = append(result, l)
	}

	return result, err
}
