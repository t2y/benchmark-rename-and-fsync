package main

import (
	"fmt"
	"io"
	"testing"
)

func createFsyncFile(t testing.TB, path string, size, syncFileRange int) {
	tmp, err := NewTempFile(path)
	if err != nil {
		t.Fatal(err)
	}

	w := NewFsyncWriter(tmp, syncFileRange)
	if _, err := io.CopyN(w, genData(size), int64(size)); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestFsyncWriter(t *testing.T) {
	dir := makeDir("fsync")
	path := fmt.Sprintf("%s/writer.txt", dir)
	createFsyncFile(t, path, 2*MiB+1, argSyncFileRange)
}

func BenchmarkFsyncWriter(b *testing.B) {
	b.ResetTimer()
	dir := ""
	for i := 0; i < b.N; i++ {
		if i%1000 == 0 {
			dir = makeDir("fsync")
		}
		path := fmt.Sprintf("%s/%04d.txt", dir, i)
		createFsyncFile(b, path, size*KiB, argSyncFileRange)
	}
}
