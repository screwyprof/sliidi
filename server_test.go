package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var (
	faultyProvider = Provider("faulty")

	clientsWithFaultyProvider = map[Provider]Client{
		faultyProvider: mockContentProvider{Source: faultyProvider, Err: errors.New("some error")},
		Provider1:      mockContentProvider{Source: Provider1},
		Provider2:      mockContentProvider{Source: Provider2},
		Provider3:      mockContentProvider{Source: Provider3},
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestAppServeHTTP(t *testing.T) {
	t.Run("it returns expected number of records", func(t *testing.T) {
		t.Parallel()

		// arrange
		want := rand.Intn(10) // nolint:gosec

		req := httptest.NewRequest("GET", "/?offset=0&count="+strconv.Itoa(want), nil)
		resp := httptest.NewRecorder()

		// act
		h := NewAppHandler(DefaultConfig, DefaultClients)
		h.ServeHTTP(resp, req)

		// assert
		assertStatusOk(t, resp.Code)
		assertResponseElementsCount(t, want, resp)
	})

	t.Run("it returns default number of records if count param is not passed", func(t *testing.T) {
		t.Parallel()

		// arrange
		want := defaultCount

		req := httptest.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()

		// act
		h := NewAppHandler(DefaultConfig, DefaultClients)
		h.ServeHTTP(resp, req)

		// assert
		assertStatusOk(t, resp.Code)
		assertResponseElementsCount(t, want, resp)
	})

	t.Run("it fallbacks to a specified provider on failure", func(t *testing.T) {
		t.Parallel()

		// arrange
		count := rand.Intn(10)                                                       // nolint:gosec
		cfg := generateConfigWithFaultyProvidersWithStableFallback(rand.Intn(9) + 1) // nolint:gosec

		want := expectedProviderQueueForConfig(cfg, count, 0)

		req := httptest.NewRequest("GET", "/?offset=0&count="+strconv.Itoa(count), nil)
		resp := httptest.NewRecorder()

		// act
		h := NewAppHandler(cfg, clientsWithFaultyProvider)
		h.ServeHTTP(resp, req)

		// assert
		assertStatusOk(t, resp.Code)
		assertConfigurationRespected(t, want, resp)
	})

	t.Run("it returns all the items before the first big failure", func(t *testing.T) {
		t.Parallel()

		// arrange
		count := rand.Intn(10)                                                                       // nolint:gosec
		cfg := generateConfigWithFaultyProvidersWithNoFallbackOrWithFaultyFallback(rand.Intn(9) + 1) // nolint:gosec

		want := expectedProviderQueueForConfig(cfg, count, 0)

		req := httptest.NewRequest("GET", "/?offset=0&count="+strconv.Itoa(count), nil)
		resp := httptest.NewRecorder()

		// act
		h := NewAppHandler(cfg, clientsWithFaultyProvider)
		h.ServeHTTP(resp, req)

		// assert
		assertStatusOk(t, resp.Code)
		assertConfigurationRespected(t, want, resp)
	})

	t.Run("it offsets the configuration by the provided number", func(t *testing.T) {
		t.Parallel()

		// arrange
		count := rand.Intn(10)  // nolint:gosec
		offset := rand.Intn(10) // nolint:gosec

		want := expectedProviderQueueForConfig(DefaultConfig, count, offset)

		url := fmt.Sprintf("/?count=%d&offset=%d", count, offset)
		req := httptest.NewRequest("GET", url, nil)
		resp := httptest.NewRecorder()

		// act
		h := NewAppHandler(DefaultConfig, DefaultClients)
		h.ServeHTTP(resp, req)

		// assert
		assertStatusOk(t, resp.Code)
		assertConfigurationRespected(t, want, resp)
	})
}

type mockContentProvider struct {
	Err    error
	Source Provider
}

func (cp mockContentProvider) GetContent(userIP string, count int) ([]*ContentItem, error) {
	if cp.Err != nil {
		return nil, cp.Err
	}

	resp := make([]*ContentItem, count)
	for i := range resp {
		resp[i] = &ContentItem{
			ID:     strconv.Itoa(rand.Int()), // nolint:gosec
			Title:  "title",
			Source: string(cp.Source),
			Expiry: time.Now(),
		}

	}
	return resp, nil
}

func generateConfigWithFaultyProvidersWithStableFallback(n int) ContentMix {
	providers := []*Provider{nil, &Provider1, &Provider2, &Provider3, &faultyProvider}

	config := make(ContentMix, 0, n)
	for i := 0; i < n; i++ {
		p := providers[rand.Intn(len(providers)-1)+1] // nolint:gosec
		fallback := selectSableFallback(providers)
		config = append(config, ContentConfig{Type: *p, Fallback: fallback})
	}

	return config
}

func selectSableFallback(providers []*Provider) *Provider {
	var fallback *Provider
	for {
		fallback = providers[rand.Intn(len(providers)-1)] // nolint:gosec
		if fallback != nil && fallback != &faultyProvider {
			break
		}
	}
	return fallback
}

func generateConfigWithFaultyProvidersWithNoFallbackOrWithFaultyFallback(n int) ContentMix {
	faultyProviders := []*Provider{nil, &faultyProvider}
	providers := []*Provider{&faultyProvider, &Provider1, &Provider2, &Provider3}

	config := make(ContentMix, 0, n)
	for i := 0; i < n; i++ {
		p := providers[rand.Intn(len(providers))]                    // nolint:gosec
		fallback := faultyProviders[rand.Intn(len(faultyProviders))] // nolint:gosec
		config = append(config, ContentConfig{Type: *p, Fallback: fallback})
	}

	return config
}

func expectedProviderQueueForConfig(cfg ContentMix, count, offset int) []Provider {
	queue := make([]Provider, 0, count)
	providersList := allStableProvidersForConfig(cfg)
	for i := 0; i < count; i++ {
		p := providersList[(i+offset)%len(providersList)]
		if p == faultyProvider {
			break
		}
		queue = append(queue, p)
	}
	return queue
}

func allStableProvidersForConfig(cfg ContentMix) []Provider {
	providers := make([]Provider, 0, len(cfg))
	for _, c := range cfg {
		if c.Type == faultyProvider && c.Fallback != nil {
			providers = append(providers, *c.Fallback)
		} else {
			providers = append(providers, c.Type)
		}
	}
	return providers
}

func assertStatusOk(t *testing.T, got int) {
	t.Helper()
	assertStatusCode(t, http.StatusOK, got)
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
	equals(t, true, ok)

	equals(t, want, len(elements))
}

func assertConfigurationRespected(t *testing.T, want []Provider, resp *httptest.ResponseRecorder) {
	t.Helper()

	var items []*ContentItem
	ok(t, json.NewDecoder(resp.Body).Decode(&items))

	got := make([]Provider, 0, len(items))
	for _, item := range items {
		got = append(got, Provider(item.Source))
	}

	assertProvidersMatch(t, want, got)
}

func assertProvidersMatch(t *testing.T, want []Provider, got []Provider) {
	t.Helper()

	equals(t, len(want), len(got))

	diff := make(map[Provider]int, len(want))
	for _, x := range want {
		diff[x]++
	}

	for _, y := range got {
		_, ok := diff[y]
		equals(t, true, ok)

		diff[y]--
		if diff[y] == 0 {
			delete(diff, y)
		}
	}

	equals(t, 0, len(diff))
}

func ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("\033[31munexpected error: %v\033[39m\n\n", err)
	}
}

func equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(exp, act) {
		tb.Fatalf("\033[31m\n\n\texp:\n\t%#+v\n\n\tgot:\n\t%#+v\033[39m", exp, act)
	}
}
