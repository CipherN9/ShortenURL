package shorten_links

import (
	"context"
	"fmt"
	"log"
	"os/exec"
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
	randNumber, err := RandString(8)
	if err != nil {
		return nil, fmt.Errorf("problem with generating short link %v", err)
	}

	shortenLink := domain + "/" + randNumber

	createdLink, err := s.Repo.Add(ctx, Link{InitialLink: initialLink, ShortenLink: shortenLink})

	if err != nil {
		return nil, err
	}

	log.Printf("New shortened link: %s", createdLink.ShortenLink)

	return createdLink, nil
}
