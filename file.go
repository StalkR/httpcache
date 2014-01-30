package httpcache

import (
	"encoding/gob"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

// A fileCache implements a cache with files saved in Path.
// It is not safe for concurrency.
type fileCache struct {
	Path string
}

// Windows does not accept `\/:*?"<>|` and UNIX `/`, replace all with dash.
var fileNameReplacer = strings.NewReplacer(`/`, `-`, `\`, `-`, `:`, `-`, `*`,
	`-`, `?`, `-`, `"`, `-`, `<`, `-`, `>`, `-`, `|`, `-`)

// fileName builds a file name from an URL to be used as cache.
func (f fileCache) fileName(u *url.URL) string {
	return path.Join(f.Path, fileNameReplacer.Replace(u.String()))
}

// Get gets data saved for an URL if present in cache.
func (f fileCache) Get(u *url.URL) (*entry, error) {

	fp, err := os.Open(f.fileName(u))
	if err != nil {
		log.Println("Could not open ", u)
		return nil, err
	}
	decoder := gob.NewDecoder(fp)
	defer fp.Close()

	var e entry
	err = decoder.Decode(&e)
	if err != nil {
		log.Printf("Could not decode %s: %s\n", u, err)
		return nil, err
	}

	return &e, nil
}

// Put puts data of an URL in cache.
func (f fileCache) Put(u *url.URL, e *entry) error {

	fp, err := os.Create(f.fileName(u))
	if err != nil {

		return err
	}
	defer fp.Close()

	encoder := gob.NewEncoder(fp)
	err = encoder.Encode(*e)
	if err != nil {
		log.Println("Could not write to cache: ", err)
		return err
	}
	return nil
}
