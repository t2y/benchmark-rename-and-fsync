package main

import (
	"io"

	"github.com/pkg/errors"
)

type FsyncWriter struct {
	t             *TempFile
	syncFileRange int64
	syncOffset    int64
	syncLen       int64
}

func NewFsyncWriter(tmp *TempFile, syncFileRange int) (w io.WriteCloser) {
	return &FsyncWriter{
		t:             tmp,
		syncFileRange: int64(syncFileRange),
	}
}

func (w *FsyncWriter) sync() (err error) {
	if err = syncFileRange(int(w.t.f.Fd()), w.syncOffset, w.syncLen, SYNC_FILE_RANGE_WRITE); err != nil {
		err = errors.Wrap(err, "call sync_file_range")
		return
	}
	w.syncOffset, w.syncLen = w.syncOffset+w.syncLen, 0
	return
}

func (w *FsyncWriter) Write(b []byte) (n int, err error) {
	n, err = w.t.Write(b)
	w.syncLen += int64(n)
	if w.syncLen >= w.syncFileRange {
		if err = w.sync(); err != nil {
			err = errors.Wrap(err, "sync in writing")
			return
		}
	}
	return
}

func (w *FsyncWriter) Close() (err error) {
	if w.syncLen > 0 {
		if err = w.sync(); err != nil {
			err = errors.Wrap(err, "sync in writing")
			return
		}
	}

	if _, err = w.t.Close(); err != nil {
		err = errors.Wrap(err, "close temp file, then renamed")
		return
	}
	return
}
