package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/i-zaitsev/hit"
)

const logo = `
 __  __     __     ______
/\ \_\ \   /\ \   /\__ __\
\ \  __ \  \ \ \  \/_/\ \/
 \ \_\ \_\  \ \_\    \ \_\
  \/_/\/_/   \/_/     \/_/

 [ the load testing tool ]`

type env struct {
	stdout io.Writer
	stderr io.Writer
	args   []string
	dryRun bool
	debug  bool
}

func main() {
	exitStatus := 0
	if err := run(&env{
		stdout: os.Stdout,
		stderr: os.Stderr,
		args:   os.Args,
	}); err != nil {
		exitStatus = 1
	}
	time.Sleep(10 * time.Millisecond)
	goroutineCheck()
	os.Exit(exitStatus)
}

func run(e *env) error {
	c := config{
		n: 100,
		c: 1,
	}
	if err := parseArgs(&c, e.args[1:], e.stderr); err != nil {
		return err
	}
	e.dryRun = c.dryRun
	e.debug = c.debug

	if e.debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	_, _ = fmt.Fprintf(
		e.stdout, "%s\n\nSending %d requests to %q (concurrency: %d)\n",
		logo, c.n, c.url, c.c,
	)
	if e.dryRun {
		return nil
	}
	if err := runHit(&c, e.stdout); err != nil {
		_, _ = fmt.Fprintf(e.stderr, "\nerror occurred: %v\n", err)
		return err
	}
	return nil
}

func runHit(c *config, stdout io.Writer) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	req, err := http.NewRequest(http.MethodGet, c.url, http.NoBody)
	if err != nil {
		return fmt.Errorf("creating a new request: %w", err)
	}

	results, err := hit.SendN(ctx, c.n, req, hit.Options{
		Concurrency: c.c,
		RPS:         c.rps,
	})

	if err != nil {
		return fmt.Errorf("sending requests: %w", err)
	}

	printSummary(hit.Summarize(results), stdout)

	return ctx.Err()
}

func printSummary(sum hit.Summary, stdout io.Writer) {
	_, _ = fmt.Fprintf(stdout, `
Summary:
	Success  : %.0f%%
	RPS	     : %.1f
	Requests : %d
	Errors	 : %d
	Bytes	 : %d
	Duration : %s
	Fastest  : %s
	Slowest  : %s
`,
		sum.Success,
		math.Round(sum.RPS),
		sum.Requests,
		sum.Errors,
		sum.Bytes,
		sum.Duration.Round(time.Millisecond),
		sum.Fastest.Round(time.Millisecond),
		sum.Slowest.Round(time.Millisecond),
	)
}

func goroutineCheck() {
	slog.Debug("runtime", "goroutines_at_exit", runtime.NumGoroutine())
	if n := runtime.NumGoroutine(); n > 1 {
		var buf strings.Builder
		_ = pprof.Lookup("goroutine").WriteTo(&buf, 1)
		for line := range strings.Lines(buf.String()) {
			slog.Debug(line)
		}
	}
}
