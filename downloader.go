package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func newDownloader(ctx context.Context) (*downloader, error) {
	dname, err := os.MkdirTemp("", "delivery-reader")
	if err != nil {
		return nil, err
	}

	return &downloader{
		tempDir: dname,
		wg:      sync.WaitGroup{},
		ctx:     ctx,
	}, nil

}

type downloader struct {
	tempDir string
	wg      sync.WaitGroup
	ctx     context.Context
}

func (d *downloader) cleanup() error {
	if err := os.RemoveAll(d.tempDir); err != nil {
		return fmt.Errorf("cleanup: %w", err)
	}
	return nil
}

func (d *downloader) download(url string) <-chan downloadResult {

	resultCh := make(chan downloadResult, 100)

	go func() {
		defer close(resultCh)

		sendResult := func(r downloadResult) {
			select {
			case <-d.ctx.Done():
			case resultCh <- r:
			}
		}

		f, err := d.fetch(url)
		if err != nil {
			sendResult(downloadResult{
				err: err,
			})
			return
		}

		sendResult(downloadResult{
			f: f,
		})

	}()

	return resultCh
}

func (d *downloader) fetch(url string) (io.ReadCloser, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: failed to download the file %w", err)
	}

	defer res.Body.Close()

	uncompressedReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("download: failed to unzip the file %w", err)
	}

	tarReader := tar.NewReader(uncompressedReader)

	f, err := os.CreateTemp(d.tempDir, "delivery-*")
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}

	tarReader.Next()
	io.Copy(f, tarReader)
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("readFile: failed to cp tar stream %w", err)
	}

	return os.Open(f.Name())
}

type downloadResult struct {
	err error
	f   io.ReadCloser
}
