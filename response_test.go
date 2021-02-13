package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponseBinder_BindResponse(t *testing.T) {
	t.Run("it returns an error status code if it cannot bind the response", func(t *testing.T) {
		rb := jsonResponseBinder{}

		req := make(chan struct{})
		resp := &httptest.ResponseRecorder{}

		rb.bindResponse(resp, req)

		equals(t, http.StatusInternalServerError, resp.Code)
	})
}
