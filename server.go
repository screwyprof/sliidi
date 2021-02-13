package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const defaultPageSize = 5
const defaultRecordsPerRequest = 1

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

	resp, err := a.fetchItems(count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.bindResponse(w, resp)
}

func (a App) fetchItems(count int) ([]*ContentItem, error) {
	resp := make([]*ContentItem, 0, count)
	for i := 0; i < count; i++ {
		items, err := a.fetchItem(i, "127.0.0.1", defaultRecordsPerRequest)
		if err != nil {
			return nil, err
		}
		resp = append(resp, items...)
	}
	return resp, nil
}

func (a App) fetchItem(n int, ip string, limit int) ([]*ContentItem, error) {
	p := a.Config[n%len(a.Config)]

	items, err := a.ContentClients[p.Type].GetContent(ip, limit)
	if err == nil {
		return items, nil
	}

	if p.Fallback == nil {
		return nil, err
	}

	return a.ContentClients[*p.Fallback].GetContent(ip, limit)
}

func (a App) bindResponse(w http.ResponseWriter, resp []*ContentItem) {
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
