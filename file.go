package httpcache

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
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
		return nil, fmt.Errorf("httpcache: could not open %v: %v", u, err)
	}
	defer fp.Close()

	// format before 2025-12-08: encoding/gob
	var e entry
	if err := gob.NewDecoder(fp).Decode(&e); err == nil {
		return &e, nil
	}
	if _, err := fp.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	// simpler format with no binary: first line is date, rest is data
	reader := bufio.NewReader(fp)
	firstLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("httpcache: could not read first line: %v", err)
	}
	saveTime, err := time.Parse(time.RFC3339, firstLine[:len(firstLine)-1])
	if err != nil {
		return nil, fmt.Errorf("httpcache: error parsing first line (date): %v", err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return &entry{data, saveTime}, nil
}

// Put puts data of an URL in cache.
func (f *fileCache) Put(u *url.URL, data []byte) error {
	fp, err := os.Create(f.fileName(u))
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err := fmt.Fprintln(fp, time.Now().Format(time.RFC3339)); err != nil {
		return err
	}
	if _, err := fp.Write(data); err != nil {
		return err
	}
	return nil
}
