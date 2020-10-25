package util

import (
	"net/http"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

type Downloader interface {
	Download(url string) (resp *http.Response, err error)
}

type CachedDownloader struct{}

func (c CachedDownloader) Download(url string) (resp *http.Response, err error) {
	cachePath, err := xdg.CacheFile("kbst/http/file")
	if err != nil {
		return resp, err
	}

	// use filepath.Dir we need a directory
	cache := diskcache.New(filepath.Dir(cachePath))

	transport := httpcache.NewTransport(cache)
	client := transport.Client()

	resp, err = client.Get(url)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
