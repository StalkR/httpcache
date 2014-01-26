package httpcache

import (
	"encoding/gob"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
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
func (f *fileCache) fileName(u *url.URL) string {
	return path.Join(f.Path, fileNameReplacer.Replace(u.String()))
}

// Get gets data saved for an URL if present in cache.
func (f *fileCache) Get(u *url.URL) (*entry, error) {
	fp, err := os.Open(f.fileName(u))
	if err != nil {
		return nil, fmt.Errorf("httpcache: could not open %u: %v", u, err)
	}
	decoder := gob.NewDecoder(fp)
	defer fp.Close()

	var e entry
	if err := decoder.Decode(&e); err != nil {
		return nil, fmt.Errorf("httpcache: could not decode %s: %v", u, err)
	}
	return &e, nil
}

// Put puts data of an URL in cache.
func (f *fileCache) Put(u *url.URL, data []byte) error {
	e := entry{data, time.Now()}
	fp, err := os.Create(f.fileName(u))
	if err != nil {
		return err
	}
	defer fp.Close()

	encoder := gob.NewEncoder(fp)
	if err := encoder.Encode(e); err != nil {
		return fmt.Errorf("httpcache: could not encode %s: %v", u, err)
	}
	return nil
}
