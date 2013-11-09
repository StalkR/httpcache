package httpcache

import (
	"io/ioutil"
	"net/url"
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
func (f fileCache) Get(u *url.URL) ([]byte, error) {
	data, err := ioutil.ReadFile(f.fileName(u))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Put puts data of an URL in cache.
func (f fileCache) Put(u *url.URL, data []byte) error {
	err := ioutil.WriteFile(f.fileName(u), data, 0644)
	if err != nil {
		return err
	}
	return nil
}
