package main

import (
	"fmt"
	"os"
)

const logo = `
 __  __     __     ______
/\ \_\ \   /\ \   /\__ __\
\ \  __ \  \ \ \  \/_/\ \/
 \ \_\ \_\  \ \_\    \ \_\
  \/_/\/_/   \/_/     \/_/

 [ the load testing tool ]`

func main() {
	c := config{
		n: 100,
		c: 1,
	}
	if err := parseArgs(&c, os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf(
		"%s\n\nSending %d requests to %q (concurrency: %d)\n",
		logo, c.n, c.url, c.c,
	)
}
