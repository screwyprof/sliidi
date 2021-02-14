package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponseBinder struct{}

func (rb jsonResponseBinder) bindResponse(w http.ResponseWriter, resp interface{}) {
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
