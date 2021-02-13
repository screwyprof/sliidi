package main

import (
	"log"
	"net/http"
	"strconv"
)

const defaultCount = 5
const defaultOffset = 0

// App represents the server's internal state.
// It holds configuration about providers and content
type App struct {
	fetcher fetcher
	rb      jsonResponseBinder
}

func NewAppHandler(cfg ContentMix, contentClients map[Provider]Client) App {
	f := fetcher{Config: cfg, ContentClients: contentClients}
	return App{fetcher: f}
}

func (a App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.String())

	count := a.count(req.URL.Query().Get("count"))
	offset := a.offset(req.URL.Query().Get("offset"))

	resp := a.fetcher.fetchItems(req.Header.Get("X-Forwarded-For"), count, offset)
	a.rb.bindResponse(w, resp)
}

func (a App) count(countParam string) int {
	count, err := strconv.Atoi(countParam)
	if err != nil {
		count = defaultCount
	}
	return count
}

func (a App) offset(offsetParam string) int {
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		offset = defaultOffset
	}
	return offset
}
