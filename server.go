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
	res, err := a.ContentClients[Provider1].GetContent("127.0.0.1", a.pageSizeFromRequest(countParam))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
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
