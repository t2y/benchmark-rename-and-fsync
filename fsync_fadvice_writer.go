package main

import (
	"io"

	"github.com/pkg/errors"
)

type FsyncFadviceWriter struct {
	t                 *TempFile
	syncFileRange     int64
	syncOffset        int64
	syncLen           int64
	flagSyncFileRange int
}

func NewFsyncFadviceWriter(tmp *TempFile, syncFileRange int, flagSyncFileRange string) (w io.WriteCloser) {
	return &FsyncFadviceWriter{
		t:                 tmp,
		syncFileRange:     int64(syncFileRange),
		flagSyncFileRange: getSyncFileRangeFlag(flagSyncFileRange),
	}
}

func (w *FsyncFadviceWriter) sync() (err error) {
	if err = syncFileRange(int(w.t.f.Fd()), w.syncOffset, w.syncLen, w.flagSyncFileRange); err != nil {
		err = errors.Wrap(err, "call sync_file_range")
		return
	}
	if err = fadvise(int(w.t.f.Fd()), w.syncOffset, w.syncLen, FADV_DONTNEED); err != nil {
		err = errors.Wrap(err, "call posix_fadvise")
		return
	}
	w.syncOffset, w.syncLen = w.syncOffset+w.syncLen, 0
	return
}

func (w *FsyncFadviceWriter) Write(b []byte) (n int, err error) {
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

func (w *FsyncFadviceWriter) Close() (err error) {
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
