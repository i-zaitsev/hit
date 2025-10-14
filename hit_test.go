package hit

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestSendStatusCode(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "/", http.NoBody)
	if err != nil {
		t.Fatalf("creating http request: %v", err)
	}

	fake := func(_ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError}, nil
	}
	client := &http.Client{Transport: roundTripFunc(fake)}
	result := Send(client, req)
	want := http.StatusInternalServerError

	if result.Status != want {
		t.Errorf("got %d, want %d", result.Status, want)
	}
}

func TestSendN(t *testing.T) {
	t.Parallel()

	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(
		func(_ http.ResponseWriter, _ *http.Request) {
			hits.Add(1)
		},
	))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL, http.NoBody)
	if err != nil {
		t.Fatalf("creating http request: %v", err)
	} else {
		nRequests := 10
		results, err := SendN(t.Context(), nRequests, req, Options{Concurrency: 5})
		if err != nil {
			t.Fatalf("SendN() err=%v, want nil", err)
		}
		for range results {
			// consume the iterator with results
		}
		if got := hits.Load(); got != int64(nRequests) {
			t.Errorf("got %d hits, want %d", got, nRequests)
		}
	}
}
