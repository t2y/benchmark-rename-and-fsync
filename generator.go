package main

import (
	"bytes"
	"crypto/rand"
	"io"
	"io/ioutil"
	mrand "math/rand"
	"time"
)

const (
	KiB = 1024
	MiB = 1024 * KiB
)

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func randomNumber(min, max int) int {
	mrand.Seed(time.Now().UnixNano())
	return mrand.Intn(max-min) + min
}

func genRandomSizes(n, min, max int) []int {
	sizes := make([]int, n)
	for i := 0; i < n; i++ {
		sizes[i] = randomNumber(min, max)
	}
	return sizes
}

func genData(size int) io.Reader {
	data := randomBytes(size)
	r := ioutil.NopCloser(bytes.NewReader(data))
	return r
}
