package main

import (
	"log"
	"net/http"
	"strconv"
)

const defaultPageSize = 5
const defaultRecordsPerRequest = 1

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

	countParam := req.URL.Query().Get("count")
	count := a.pageSizeFromRequest(countParam)

	resp := a.fetcher.fetchItems(req.Header.Get("X-Forwarded-For"), count, defaultRecordsPerRequest)
	a.rb.bindResponse(w, resp)
}

func (a App) pageSizeFromRequest(countParam string) int {
	count, err := strconv.Atoi(countParam)
	if err != nil {
		count = defaultPageSize
	}
	return count
}
