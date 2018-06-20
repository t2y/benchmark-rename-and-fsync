package main

import (
	"fmt"
	"log"
	"os"
)

func removeFile(path string) {
	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}
}

func makeDir(sub string) (path string) {
	path = fmt.Sprintf("%s/%s", testDir, sub)
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatal(err)
	}
	return
}
