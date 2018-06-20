package main

import (
	"fmt"
	"testing"
)

func TestFsyncFadviceWriter(t *testing.T) {
	dir := makeDir("fsync+fadvice")
	path := fmt.Sprintf("%s/writer.txt", dir)
	createFsyncFadviceFile(path, 2*MiB+1, argSyncFileRange, flagSyncFileRange)
}
