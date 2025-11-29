package main

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
