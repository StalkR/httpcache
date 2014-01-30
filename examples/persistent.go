// Binary persistent demonstrates using a persistent (with files) httpcache.
// Run the binary twice with the same URL to use the cache.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/EverythingMe/httpcache"
)

var cache = flag.String("cache", "", "directory to use as cache")

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-cache <dir>] <url>\n", os.Args[0])
		os.Exit(1)
	}
	url := flag.Arg(0)

	client := httpcache.NewPersistentClient(*cache, httpcache.NeverExpirePolicy{})
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(2)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(3)
	}
	fmt.Print(string(body))
}
