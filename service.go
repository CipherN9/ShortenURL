package main

import (
	"context"
	"fmt"
	"log"
)

type ILinksService interface {
	ResolveLink(context.Context, string) (*Link, error)
	GetLinks(context.Context, *Filter) ([]Link, error)
	ShortenLink(context.Context, string, string) (string, error)
}
type LinksService struct {
	repo *LinksRepository
}

func (s *LinksService) ResolveLink(ctx context.Context, url string) (*Link, error) {
	resolvedLinks, err := s.GetLinks(ctx, &Filter{ShortenLink: url})
	if err != nil {
		return nil, err
	}

	return &resolvedLinks[0], nil
}

func (s *LinksService) GetLinks(ctx context.Context, f *Filter) ([]Link, error) {
	links, err := s.repo.Get(ctx, f)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (s *LinksService) ShortenLink(ctx context.Context, domain string, initialLink string) (string, error) {
	randNumber, err := RandString(8)
	if err != nil {
		return "", fmt.Errorf("problem with generating short link %v", err)
	}

	shortenLink := domain + "/" + randNumber

	err = s.repo.Add(ctx, &Link{InitialLink: initialLink, ShortenLink: shortenLink})

	if err != nil {
		return "", err
	}

	log.Printf("New shortened link: %s", shortenLink)

	return shortenLink, nil
}
