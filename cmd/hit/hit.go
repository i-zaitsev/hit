package main

import (
	"fmt"
	"io"
	"os"
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
}

func main() {
	if err := run(&env{
		stdout: os.Stdout,
		stderr: os.Stderr,
		args:   os.Args,
	}); err != nil {
		os.Exit(1)
	}
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
	return nil
}
