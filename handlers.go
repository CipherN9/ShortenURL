package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type DB struct {
	Links map[string]string
}

var db = DB{Links: map[string]string{}}

func HandleResolveLink(w http.ResponseWriter, r *http.Request) {
	urlLink := ResolveDomain(r) + r.URL.Path
	log.Printf("Received request for resolve link %s \n", urlLink)

	resolvedLink, ok := db.Links[urlLink]

	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, resolvedLink, http.StatusTemporaryRedirect)
}

func HandleGetLinks(w http.ResponseWriter, _ *http.Request) {
	log.Println("Handle Get Links")

	w.Header().Set("Content-Type", "application/json")

	var linksResponse []GetLinksResponse

	for shortenLink, initialLink := range db.Links {
		linksResponse = append(linksResponse, GetLinksResponse{ShortenLink: shortenLink, InitialLink: initialLink})
	}

	log.Printf("Links response: %+v", linksResponse)

	if err := json.NewEncoder(w).Encode(&linksResponse); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func HandleLinkShorten(w http.ResponseWriter, r *http.Request) {
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
	db.Links[newLink] = l.Link

	log.Printf("New shortened link: %s", newLink)

	if err := json.NewEncoder(w).Encode(&PostLinkResponse{Link: newLink}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
