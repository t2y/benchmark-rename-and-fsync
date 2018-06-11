package main

import (
	"fmt"
	"io"
	"testing"
)

func createFsyncFadviceFile(t testing.TB, path string, size, syncFileRange int) {
	tmp, err := NewTempFile(path)
	if err != nil {
		t.Fatal(err)
	}

	w := NewFsyncFadviceWriter(tmp, syncFileRange)
	if _, err := io.CopyN(w, genData(size), int64(size)); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestFsyncFadviceWriter(t *testing.T) {
	dir := makeDir("fsync+fadvice")
	path := fmt.Sprintf("%s/writer.txt", dir)
	createFsyncFadviceFile(t, path, 2*MiB+1, argSyncFileRange)
}

func BenchmarkFsyncFadviceWriter(b *testing.B) {
	b.ResetTimer()
	dir := ""
	for i := 0; i < b.N; i++ {
		if i%1000 == 0 {
			dir = makeDir("fsync+fadvice")
		}
		path := fmt.Sprintf("%s/%04d.txt", dir, i)
		createFsyncFadviceFile(b, path, size*KiB, argSyncFileRange)
	}
}