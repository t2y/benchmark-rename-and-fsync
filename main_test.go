package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	defaultSyncFileRange = 1048576 // 1MiB
)

var (
	size             int
	testDir          string
	argSyncFileRange int

	seq *Sequence
)

type Sequence struct {
	m       sync.Mutex
	current int
}

func (s *Sequence) pop() (v int) {
	s.m.Lock()
	defer s.m.Unlock()

	v = s.current
	s.current += 1
	return
}

func init() {
	flag.IntVar(&size, "size", 1, "size of writing file (KB)")
	flag.StringVar(&testDir, "testDir", "testdata", "test data directory")
	flag.IntVar(&argSyncFileRange, "syncFileRange", defaultSyncFileRange, "size of sync_file_range(B)")

	seq = &Sequence{}
}

func removeFile(path string) {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}
}

func removeDir() {
	if err := os.RemoveAll(testDir); err != nil {
		log.Fatal(err)
	}
}

func makeDir(prefix string) (path string) {
	n := seq.pop()
	path = fmt.Sprintf("%s/%s-sub%04d", testDir, prefix, n)
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}
	return
}
