package hit

import (
	"log/slog"
	"net/http"
	"time"
)

type OptLevel int

const (
	OptDefault OptLevel = iota
	OptDisabled
	OptEnabled
)

// SendFunc is a type of function that sends an [net/http.Request] and returns a [Result].
type SendFunc func(r *http.Request) Result

// Options defines the options for sending requests.
// Use default options for unset options.
type Options struct {
	// Concurrency is the number of concurrent requests to send.
	// Default: 1
	Concurrency int

	// RPS is the requests to send per second.
	// Default: 0 (no rate limiting)
	RPS int

	// Send processes requests.
	Send SendFunc

	// Optimized replaces the [net/http.DefaultClient] with an optimized configuration if requested.
	Optimized OptLevel
}

// Defaults returns the default [Options].
func Defaults() Options {
	return withDefaults(Options{})
}

func withDefaults(o Options) Options {
	if o.Concurrency == 0 {
		slog.Warn("zero concurrency requested: running sequentially")
		o.Concurrency = 1
	}
	if o.Send == nil {
		slog.Debug("no Send provided: using the default impl")
		var client *http.Client
		switch o.Optimized {
		case OptDefault, OptEnabled:
			slog.Debug("optimized client requests: setting up transport and disabling redirects")
			client = &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost: o.Concurrency,
				},
				CheckRedirect: disableHttpRedirects,
				Timeout:       30 * time.Second,
			}
		case OptDisabled:
			slog.Debug("using default http client")
			client = http.DefaultClient
		}
		o.Send = func(r *http.Request) Result {
			return Send(client, r)
		}
	}
	return o
}

func disableHttpRedirects(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}
