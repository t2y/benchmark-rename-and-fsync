package main

import (
	"context"
	"io"
	"log"

	"github.com/pkg/errors"
)

type NoSyncWriter struct {
	t *TempFile
}

func NewNoSyncWriter(tmp *TempFile) (w io.WriteCloser) {
	return &NoSyncWriter{
		t: tmp,
	}
}

func (w *NoSyncWriter) Write(b []byte) (n int, err error) {
	n, err = w.t.Write(b)
	return
}

func (w *NoSyncWriter) Close() (err error) {
	if _, e := w.t.Close(); e != nil {
		err = errors.Wrap(e, "close temp file, then renamed")
		return
	}
	return
}

func createNoSyncFile(path string, size int) {
	tmp, err := NewTempFile(path)
	if err != nil {
		log.Fatal(err)
	}

	w := NewNoSyncWriter(tmp)
	if _, err := io.CopyN(w, genData(size), int64(size)); err != nil {
		log.Fatal(err)
	}

	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
}

func runBenchmarkNoSyncWriter(ctx context.Context, pathCh chan string, n int) (i int) {
	for {
		select {
		case <-ctx.Done():
			return // expect timeout
		default:
			path, ok := <-pathCh
			if !ok {
				return
			}
			createNoSyncFile(path, size*KiB)
			i += 1
		}
	}
}
