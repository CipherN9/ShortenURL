package short_links

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/jackc/pgx/v5"
)

type ILinksService interface {
	ResolveLink(context.Context, string) (*Link, error)
	GetLinks(context.Context, *Filter) ([]Link, error)
	ShortenLink(context.Context, string, string) (*Link, error)
}
type LinksService struct {
	Repo ILinksRepository
}

func (s *LinksService) ResolveLink(ctx context.Context, url string) (*Link, error) {
	resolvedLinks, err := s.GetLinks(ctx, &Filter{ShortenLink: url})
	if err != nil {
		return nil, err
	}

	if len(resolvedLinks) == 0 {
		return nil, exec.ErrNotFound
	}

	return &resolvedLinks[0], nil
}

func (s *LinksService) GetLinks(ctx context.Context, f *Filter) ([]Link, error) {
	links, err := s.Repo.Get(ctx, f)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (s *LinksService) ShortenLink(ctx context.Context, initialLink string, domain string) (*Link, error) {
	maxRetries := 5
	var createdLink *Link

	for i := 0; i < maxRetries; i++ {
		slug, err := RandSlug(8)
		if err != nil {
			return nil, fmt.Errorf("problem with generating short link %v", err)
		}

		shortenLink := domain + "/" + slug

		createdLink, err = s.Repo.Add(ctx, Link{InitialLink: initialLink, ShortenLink: shortenLink})

		// In case there is collision during RandString slug generation.
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	log.Printf("New shortened link: %s", createdLink.ShortenLink)

	return createdLink, nil
}
