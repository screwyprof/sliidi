package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func init() {
	rand.Seed(0)
}

func TestAppServeHTTP(t *testing.T) {
	t.Run("it returns expected number of records", func(t *testing.T) {
		t.Parallel()

		// arrange
		want := rand.Intn(10) // nolint:gosec

		req := httptest.NewRequest("GET", "/?offset=0&count="+strconv.Itoa(want), nil)
		resp := httptest.NewRecorder()

		// act
		h := newAppHandler()
		h.ServeHTTP(resp, req)

		// assert
		assertStatusCode(t, http.StatusOK, resp.Code)
		assertResponseElementsCount(t, want, resp)
	})

	t.Run("it returns default number of records if count param is not passed", func(t *testing.T) {
		t.Parallel()

		// arrange
		want := defaultPageSize

		req := httptest.NewRequest("GET", "/?offset=0", nil)
		resp := httptest.NewRecorder()

		// act
		h := newAppHandler()
		h.ServeHTTP(resp, req)

		// assert
		assertStatusCode(t, http.StatusOK, resp.Code)
		assertResponseElementsCount(t, want, resp)
	})

	t.Run("it fetches records from the given providers according to the configuration", func(t *testing.T) {
		t.Parallel()

		// arrange
		count := 8
		want := []Provider{Provider1, Provider1, Provider2, Provider3, Provider1, Provider1, Provider1, Provider2}

		req := httptest.NewRequest("GET", "/?offset=0&count="+strconv.Itoa(count), nil)
		resp := httptest.NewRecorder()

		// act
		h := newAppHandler()
		h.ServeHTTP(resp, req)

		// assert
		assertStatusCode(t, http.StatusOK, resp.Code)
		assertConfigurationRespected(t, want, resp)
	})
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

func assertStatusCode(t *testing.T, want, got int) {
	t.Helper()
	if want != got {
		t.Fatalf("want response code %d, got %d", want, got)
	}
}

func assertResponseElementsCount(t *testing.T, want int, resp *httptest.ResponseRecorder) {
	t.Helper()

	var res interface{}
	err := json.NewDecoder(resp.Body).Decode(&res)
	ok(t, err)

	elements, ok := res.([]interface{})
	true(t, ok)

	equals(t, want, len(elements))
}

func assertConfigurationRespected(t *testing.T, want []Provider, resp *httptest.ResponseRecorder) {
	var items []*ContentItem
	ok(t, json.NewDecoder(resp.Body).Decode(&items))

	got := make([]Provider, 0, len(items))
	for _, item := range items {
		got = append(got, Provider(item.Source))
	}

	equals(t, want, got)
}

func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("\033[31munexpected error: %v\033[39m\n\n", err)
	}
}

func true(tb testing.TB, condition bool) {
	tb.Helper()
	if !condition {
		tb.Errorf("\033[31mcondition is false\033[39m\n\n")
	}
}

func equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(exp, act) {
		tb.Errorf("\033[31m\n\n\texp:\n\t%#+v\n\n\tgot:\n\t%#+v\033[39m", exp, act)
	}
}
