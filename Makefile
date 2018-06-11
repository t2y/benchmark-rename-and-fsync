GOPATH := $(shell pwd)
export GOPATH

UNITTEST_DIR := testdata-single
BENCHTEST_DIR := testdata

all: build

deps:
	go get -v -d .

build:
	go build -o main .

test:
	rm -rf $(UNITTEST_DIR)
	go test . -v -testDir $(UNITTEST_DIR) -syncFileRange 102400

bench:
	rm -rf $(BENCHTEST_DIR)
	go test -v -benchmem -bench . -testDir $(BENCHTEST_DIR) -count 3 -cpu 4 -benchtime 1s -size 4 -run Benchmark*

1KB:
	go test -bench . -v -benchmem -run Benchmark* -benchtime 5s -size 1

100KB:
	go test -bench . -v -benchmem -run Benchmark* -benchtime 5s -size 100

1MB:
	go test -bench . -v -benchmem -run Benchmark* -benchtime 5s -size 1024

1MBover:
	go test -bench . -v -benchmem -run Benchmark* -benchtime 5s -size 1025

clean:
	rm -f main
