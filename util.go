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
	for {
		dirNum = randomNumber(0, m.maxDirNum)
		fileNum = m.dirs[dirNum]
		if fileNum < 1000 {
			m.dirs[dirNum] = fileNum + 1
			break
		}
	}
	path = fmt.Sprintf("%s/%s-%07d/%03d.txt", testDir, m.prefix, dirNum, fileNum)
	return
}

func NewRandomDirMaker(prefix string) (m *RandomDirMaker) {
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
	for i := 0; i < maxRandomDirNum; i++ {
		path := fmt.Sprintf("%s/%s-%07d", testDir, prefix, i)
		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatal(err)
		}
	}
}
