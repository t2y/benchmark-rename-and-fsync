package main

import (
	"fmt"
	"log"
	"math"
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

type directoryMaker interface {
	create() (path string)
}

type SequentialDirMaker struct {
	prefix string
	dir    string
	dirNum uint64
	i      uint64
}

func (m *SequentialDirMaker) create() (path string) {
	if m.i%1000 == 0 {
		m.dir = makeDir(fmt.Sprintf("%s-%05d", m.prefix, m.dirNum))
		m.dirNum += 1
		m.i = 0
	}
	path = fmt.Sprintf("%s/%03d.txt", m.dir, m.i)
	m.i += 1
	return
}

type RandomDirMaker struct {
	prefix    string
	maxDirNum int
	dirs      map[int]int
}

func (m *RandomDirMaker) create() (path string) {
	var (
		dirNum  int
		fileNum int
	)
	dirNum = randomNumber(0, m.maxDirNum)
	fileNum = m.dirs[dirNum]
	m.dirs[dirNum] = fileNum + 1
	path = fmt.Sprintf("%s/%s-%04d/%03d.txt", testDir, m.prefix, dirNum, fileNum)
	return
}

func getMaxRandomDir() (n int) {
	// 1 directory has about 1000 files
	v := float64(defaultDiskThroughput*duration.Seconds()) / float64(size*KiB) / 1000.0
	n = int(math.Floor(v + 0.5))
	if n == 0 {
		n = 1
	}
	return
}

func NewRandomDirMaker(prefix string) (m *RandomDirMaker) {
	maxRandomDirNum := getMaxRandomDir()
	log.Printf("number of random directories: %d", maxRandomDirNum)
	dirs := make(map[int]int, maxRandomDirNum)
	for i := 0; i < maxRandomDirNum; i++ {
		dirs[i] = 0
	}
	m = &RandomDirMaker{
		prefix:    prefix,
		dirs:      dirs,
		maxDirNum: maxRandomDirNum,
	}
	return
}

func createRandomDirectory(prefix string) {
	maxRandomDirNum := getMaxRandomDir()
	for i := 0; i < maxRandomDirNum; i++ {
		path := fmt.Sprintf("%s/%s-%04d", testDir, prefix, i)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatal(err)
		}
	}
}
