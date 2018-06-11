package main

import (
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	SYNC_FILE_RANGE_WRITE = 2

	FADV_DONTNEED = 0x4
)

func syncFileRange(fd int, off int64, n int64, flags int) error {
	return syscall.SyncFileRange(fd, off, n, flags)
}

func fadvise(fd int, offset int64, length int64, advice int) error {
	return unix.Fadvise(fd, offset, length, advice)
}
