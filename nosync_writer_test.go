package main

import (
	"fmt"
	"io"
	"testing"
)

func createNoSyncFile(t testing.TB, path string, size int) {
	tmp, err := NewTempFile(path)
	if err != nil {
		t.Fatal(err)
	}

	w := NewNoSyncWriter(tmp)
	if _, err := io.CopyN(w, genData(size), int64(size)); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNoSyncWriter(t *testing.T) {
	dir := makeDir("nosync")
	path := fmt.Sprintf("%s/writer.txt", dir)
	createNoSyncFile(t, path, 1*KiB)
}

func BenchmarkNoSyncWriter(b *testing.B) {
	b.ResetTimer()
	dir := ""
	for i := 0; i < b.N; i++ {
		if i%1000 == 0 {
			dir = makeDir("nosync")
		}
		path := fmt.Sprintf("%s/%04d.txt", dir, i)
		createNoSyncFile(b, path, size*KiB)
	}
}
