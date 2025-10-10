package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type config struct {
	url string
	n   int
	c   int
	rps int
}

func parseArgs(c *config, args []string) error {
	fs := flag.NewFlagSet("hit", flag.ContinueOnError)
	fs.StringVar(&c.url, "url", "", "HTTP server `URL` (required)")
	fs.IntVar(&c.n, "n", c.n, "Number of requests")
	fs.IntVar(&c.c, "c", c.c, "Concurrency level")
	fs.IntVar(&c.rps, "rps", c.rps, "Requests per second")
	return fs.Parse(args)
}

func parseArgsCustom(c *config, args []string) error {
	flagSet := map[string]parseFunc{
		"url": stringVar(&c.url),
		"n":   intVar(&c.n),
		"c":   intVar(&c.c),
		"rps": intVar(&c.rps),
	}
	for _, arg := range args {
		name, val, _ := strings.Cut(arg, "=")
		name = strings.TrimPrefix(name, "-")
		setVar, ok := flagSet[name]
		if !ok {
			return fmt.Errorf("flag provided but not defined: -%s", name)
		}
		if err := setVar(val); err != nil {
			return fmt.Errorf("invalid value %q for flag -%s: %w",
				val, name, err,
			)
		}
	}
	return nil
}

type parseFunc func(string) error

func stringVar(p *string) parseFunc {
	return func(s string) error {
		*p = s
		return nil
	}
}

func intVar(p *int) parseFunc {
	return func(s string) error {
		var err error
		*p, err = strconv.Atoi(s)
		return err
	}
}
