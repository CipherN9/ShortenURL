package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleResolveLink(repo *LinksRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		urlLink := ResolveDomain(r) + r.URL.Path
		log.Printf("Received request for resolve link %s \n", urlLink)

		ctx := r.Context()

		resolvedLinks, err := repo.Get(ctx, &Link{ShortenLink: urlLink})

		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, resolvedLinks[0].InitialLink, http.StatusTemporaryRedirect)
	}
}

func HandleGetLinks(repo *LinksRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handle Get Links")

		w.Header().Set("Content-Type", "application/json")

		var linksResponse []GetLinksResponse

		ctx := r.Context()

		links, err := repo.Get(ctx, &Link{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		for _, link := range links {
			linksResponse = append(linksResponse, GetLinksResponse{ShortenLink: link.ShortenLink, InitialLink: link.InitialLink})
		}

		log.Printf("Links response: %+v", linksResponse)

		if err := json.NewEncoder(w).Encode(&linksResponse); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func HandleLinkShorten(repo *LinksRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		randNumber, err := RandString(8)
		if err != nil {
			http.Error(w, "Problem with generating short link", http.StatusBadRequest)
		}

		newLink := ResolveDomain(r) + "/" + randNumber

		ctx := r.Context()
		err = repo.Add(ctx, &Link{InitialLink: l.Link, ShortenLink: newLink})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		log.Printf("New shortened link: %s", newLink)

		if err := json.NewEncoder(w).Encode(&PostLinkResponse{Link: newLink}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
