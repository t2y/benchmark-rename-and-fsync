package main

import (
	"io"

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
