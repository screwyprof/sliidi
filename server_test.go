package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppServeHTTP(t *testing.T) {
	// arrange
	h := App{
		ContentClients: map[Provider]Client{
			Provider1: SampleContentProvider{Source: Provider1},
			Provider2: SampleContentProvider{Source: Provider2},
			Provider3: SampleContentProvider{Source: Provider3},
		},
		Config: DefaultConfig,
	}

	s := httptest.NewServer(h)
	defer s.Close()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/?offset=0&count=5", nil)

	// act
	h.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusOK {
		t.Fatalf("Response code is %d, want 200", rr.Code)
	}
}
