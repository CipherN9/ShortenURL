package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleResolveLink(service ILinksService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		urlLink := ResolveDomain(r) + r.URL.Path
		log.Printf("Received request for resolve link %s \n", urlLink)

		resolved, err := service.ResolveLink(ctx, urlLink)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, resolved.InitialLink, http.StatusTemporaryRedirect)
	}
}

func HandleGetLinks(service ILinksService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Println("Handle Get Links")

		w.Header().Set("Content-Type", "application/json")

		links, err := service.GetLinks(ctx, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var linksResponse []GetLinksResponse
		for _, link := range links {
			linksResponse = append(linksResponse, GetLinksResponse{ShortenLink: link.ShortenLink, InitialLink: link.InitialLink})
		}

		log.Printf("Links response: %+v", linksResponse)

		if err := json.NewEncoder(w).Encode(linksResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HandleLinkShorten(service ILinksService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Content-Type", "application/json")

		var l PostLinkPayload
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Handle Link Shorten for: %v", l.Link)

		if !IsValidURL(l.Link) {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		links, err := service.GetLinks(ctx, &Filter{InitialLink: l.Link})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(links) == 1 {
			if err := json.NewEncoder(w).Encode(PostLinkResponse{Link: links[0].ShortenLink}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		shortenLink, err := service.ShortenLink(ctx, l.Link, ResolveDomain(r))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(PostLinkResponse{Link: shortenLink}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
