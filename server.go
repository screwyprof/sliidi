package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const defaultPageSize = 5

// App represents the server's internal state.
// It holds configuration about providers and content
type App struct {
	ContentClients map[Provider]Client
	Config         ContentMix
}

func (a App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.String())

	countParam := req.URL.Query().Get("count")
	count := a.pageSizeFromRequest(countParam)

	// fetch items
	resp := make([]*ContentItem, 0, count)
	for i := 0; i < count; i++ {
		items, err := a.ContentClients[a.Config[i%len(a.Config)].Type].GetContent("127.0.0.1", 1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp = append(resp, items...)
	}

	// marshal response
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a App) pageSizeFromRequest(countParam string) int {
	count, err := strconv.Atoi(countParam)
	if err != nil {
		count = defaultPageSize
	}
	return count
}
