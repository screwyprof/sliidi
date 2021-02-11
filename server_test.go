package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppServeHTTP(t *testing.T) {
	// arrange
	h := newAppHandler()

	req := httptest.NewRequest("GET", "/?offset=0&count=5", nil)
	resp := httptest.NewRecorder()

	// act
	h.ServeHTTP(resp, req)

	// assert
	assertStatusCode(t, http.StatusOK, resp.Code)
}

func assertStatusCode(t *testing.T, want, got int) {
	t.Helper()
	if want != got {
		t.Fatalf("want response code %d, got %d", want, got)
	}
}

func newAppHandler() App {
	h := App{
		ContentClients: map[Provider]Client{
			Provider1: SampleContentProvider{Source: Provider1},
			Provider2: SampleContentProvider{Source: Provider2},
			Provider3: SampleContentProvider{Source: Provider3},
		},
		Config: DefaultConfig,
	}
	return h
}
