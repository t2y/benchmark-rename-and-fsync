// +build !linux

package main

const (
	SYNC_FILE_RANGE_WRITE = 2

	FADV_DONTNEED = 0x4
)

func syncFileRange(fd int, off int64, n int64, flags int) error { return nil }

func fadvise(fd int, offset int64, length int64, advice int) error { return nil }
