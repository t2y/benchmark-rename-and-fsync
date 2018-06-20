package main

import (
	"fmt"
	"testing"
)

func TestNoSyncWriter(t *testing.T) {
	dir := makeDir("nosync")
	path := fmt.Sprintf("%s/writer.txt", dir)
	createNoSyncFile(path, 1*KiB)
}
