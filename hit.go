package hit

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func Send(_ *http.Client, _ *http.Request) Result {
	// placeholder: real impl comes in the next changes

	const roundTripTime = 100 * time.Millisecond

	time.Sleep(roundTripTime)

	return Result{
		Status:   http.StatusOK,
		Bytes:    10,
		Duration: roundTripTime,
	}
}

// SendN sends N requests using [Send].
// It returns a single-user [Results] iterator that pushes a [Result] for each [net/http.Request] sent.
func SendN(
	ctx context.Context,
	n int, req *http.Request, opts Options,
) (Results, error) {
	if n <= 0 {
		return nil, fmt.Errorf("n must be positive: got %d", n)
	}

	ctx, cancel := context.WithCancel(ctx)
	results := runPipeline(ctx, n, req, withDefaults(opts))

	return func(yield func(Result) bool) {
		defer cancel()
		for result := range results {
			if !yield(result) {
				return
			}
		}
	}, nil
}
