// Binary persistent demonstrates using a persistent (with files) httpcache.
// Run the binary twice with the same URL to use the cache.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/StalkR/httpcache"
)

var (
	cache = flag.String("cache", ".", "directory to use as cache")
	ttl   = flag.Duration("ttl", time.Minute, "cache expiration")
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-cache <dir>] [-ttl <duration>] <url>\n", os.Args[0])
		os.Exit(1)
	}
	url := flag.Arg(0)

	client, err := httpcache.NewPersistentClient(*cache, *ttl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(3)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(4)
	}
	fmt.Print(string(body))
}
