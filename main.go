package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type DB struct {
	Links map[string]string
}

type PostLinkPayload struct {
	Link string `json:"link"`
}

type GetLinksResponse struct {
	InitialLink string `json:"initialLink"`
	ShortenLink string `json:"shortenLink"`
}

type PostLinkResponse struct {
	Link string `json:"link"`
}

var db = DB{Links: map[string]string{}}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func IsValidURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func RandString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i := 0; i < n; i++ {
		b[i] = letters[int(b[i])%len(letters)]
	}

	return string(b), nil
}

func ResolveDomain(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + r.Host
}

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

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /{shorten}", HandleResolveLink)
	router.HandleFunc("GET /links", HandleGetLinks)
	router.HandleFunc("POST /links/shorten", HandleLinkShorten)

	server := http.Server{Addr: ":8080", Handler: router}
	log.Println("Start server on port 8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server stopped running.", err)
	}
}
