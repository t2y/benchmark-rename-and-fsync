package main

import (
	"io"

	"github.com/pkg/errors"
)

type FsyncWriter struct {
	t                 *TempFile
	syncFileRange     int64
	syncOffset        int64
	syncLen           int64
	flagSyncFileRange int
}

func NewFsyncWriter(tmp *TempFile, syncFileRange int, flagSyncFileRange string) (w io.WriteCloser) {
	return &FsyncWriter{
		t:                 tmp,
		syncFileRange:     int64(syncFileRange),
		flagSyncFileRange: getSyncFileRangeFlag(flagSyncFileRange),
	}
}

func (w *FsyncWriter) sync() (err error) {
	if err = syncFileRange(int(w.t.f.Fd()), w.syncOffset, w.syncLen, w.flagSyncFileRange); err != nil {
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

func getSyncFileRangeFlag(argFlag string) (flag int) {
	switch argFlag {
	case "before":
		flag = SYNC_FILE_RANGE_WAIT_BEFORE
	case "write":
		flag = SYNC_FILE_RANGE_WRITE
	case "after":
		flag = SYNC_FILE_RANGE_WAIT_AFTER
	case "before-write":
		flag = SYNC_FILE_RANGE_WAIT_BEFORE + SYNC_FILE_RANGE_WRITE
	case "before-write-after":
		flag = SYNC_FILE_RANGE_WAIT_BEFORE + SYNC_FILE_RANGE_WRITE + SYNC_FILE_RANGE_WAIT_AFTER
	default:
		flag = SYNC_FILE_RANGE_WRITE
	}
	return
}
