package httpcache

import (
  "bytes"
  "net/url"
  "testing"
)

func TestFileCache(t *testing.T) {
  f := fileCache{Path: t.TempDir()}
  u, err := url.Parse("https://example.com/pa/th?q=v")
  if err != nil {
    t.Fatal(err)
  }
  data := []byte("some\ndata")
  if err := f.Put(u, data); err != nil {
    t.Fatal(err)
  }
  e, err := f.Get(u)
  if err != nil {
    t.Fatal(err)
  }
  if e.SaveTime.IsZero() {
    t.Errorf("save time not set")
  }
  if !bytes.Equal(data, e.Data) {
    t.Errorf("wrong data: got %v, want %v", string(e.Data), string(data))
  }
}
