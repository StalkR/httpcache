// Binary volatile demonstrates using a volatile (in memory) httpcache.
// Run the binary once with repeated arguments to use the cache.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/EverythingMe/httpcache"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s <url> [<url> ...]\n", os.Args[0])
		os.Exit(1)
	}

	client := httpcache.NewVolatileClient(httpcache.NeverExpirePolicy{}, 1000)
	for _, url := range flag.Args() {
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
}
