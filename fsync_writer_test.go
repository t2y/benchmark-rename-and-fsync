package main

import (
	"fmt"
	"testing"
)

func TestFsyncWriter(t *testing.T) {
	dir := makeDir("fsync")
	path := fmt.Sprintf("%s/writer.txt", dir)
	createFsyncFile(path, 2*MiB+1, argSyncFileRange, flagSyncFileRange)
}
