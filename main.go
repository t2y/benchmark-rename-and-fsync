package main

import (
	"context"
	"flag"
	"log"
	"math"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof"
)

const (
	defaultSyncFileRange  = 1048576 // 1MiB
	defaultDiskThroughput = 5 * MiB // 5MiB/sec
)

var (
	size              int
	testDir           string
	argSyncFileRange  int
	flagSyncFileRange string
	benchmark         string
	dirMaker          string

	server      bool
	concurrent  int
	duration    time.Duration
	enablePprof bool
	verbose     bool
)

type benchmarkFunc func(context.Context, chan string, int) int

func getBenchmark() (f benchmarkFunc) {
	switch benchmark {
	case "nosync":
		return runBenchmarkNoSyncWriter
	case "fsync":
		return runBenchmarkFsyncWriter
	case "fsyn+fadv":
		return runBenchmarkFsyncFadviceWriter
	default:
		log.Fatalf("unknown benchmark function: %s", benchmark)
	}
	return
}

func initFlags() {
	flag.IntVar(&size, "size", 1, "size of writing file (KB)")
	flag.StringVar(&testDir, "testDir", "testdata", "test data directory")
	flag.StringVar(&flagSyncFileRange, "syncFileRangeFlag", "write", "flag for sync_file_range")
	flag.IntVar(&argSyncFileRange, "syncFileRange", defaultSyncFileRange, "size of sync_file_range(B)")
	flag.StringVar(&benchmark, "benchmark", "", "choose nosync|fsync|fsyn+fadv")
	flag.StringVar(&dirMaker, "dirMaker", "", "choose sequential|random")

	flag.BoolVar(&server, "server", false, "set server mode")
	flag.IntVar(&concurrent, "concurrent", 2, "number of goroutines")
	flag.DurationVar(&duration, "duration", 3*time.Second, "run benchmark (e.g. 10s, 1m)")
	flag.BoolVar(&enablePprof, "pprof", false, "enable pprof")
	flag.BoolVar(&verbose, "verbose", false, "set verbose mode")
}

func getPathChan(dirMaker directoryMaker) (pathCh chan string) {
	pathCh = make(chan string, 256)
	go func(dirMaker directoryMaker) {
		for {
			pathCh <- dirMaker.create()
		}
	}(dirMaker)
	return
}

func runLocalBench(pathCh chan string) {
	resultCh := make(chan int, concurrent)
	benchmarkFunc := getBenchmark()

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	startTime := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < int(concurrent); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			numberOfCreated := benchmarkFunc(ctx, pathCh, i)
			if verbose {
				log.Printf("goroutine %d, number of created: %d\n", i, numberOfCreated)
			}
			resultCh <- numberOfCreated
		}(i)
	}
	wg.Wait()

	elapsedTime := time.Since(startTime)
	log.Printf("time duration is %v, it took %s", duration, elapsedTime)

	total := 0
	for i := 0; i < concurrent; i++ {
		total += <-resultCh
	}
	close(resultCh)

	filesPerSec := float64(total) / elapsedTime.Seconds()
	nanoSecPerFile := elapsedTime.Nanoseconds() / int64(total)
	totalQuantity := float64(total*size*KiB) / 1024 / 1024
	throughput := totalQuantity / elapsedTime.Seconds()

	log.Printf("total number of created: %d", total)
	log.Printf("- files/second: %v", math.Floor(filesPerSec+0.5))
	log.Printf("- nanoseconds/file: %v", nanoSecPerFile)
	log.Printf("- total quantity (MiB): %v", totalQuantity)
	log.Printf("- throughput (MiB/sec): %v", throughput)
}

func runServerBench(pathCh chan string) {
	http.Handle("/fsync", &fsyncFadviceHandler{
		size:   size,
		pathCh: pathCh,
	})
	http.Handle("/nosync", &nosyncHandler{
		size:   size,
		pathCh: pathCh,
	})

	addr := "localhost:8090"
	log.Printf("server start at %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	initFlags()
	flag.Parse()

	if enablePprof {
		go func() {
			addr := "localhost:9090"
			log.Printf("enable pprof at %s\n", addr)
			log.Println(http.ListenAndServe(addr, nil))
		}()
	}

	length := 256
	if concurrent > length {
		length = concurrent
	}
	pathCh := make(chan string, length)

	maker := getDirMaker(dirMaker)
	go func(maker directoryMaker) {
		for {
			pathCh <- maker.create()
		}
	}(maker)

	if server {
		runServerBench(pathCh)
	} else {
		runLocalBench(pathCh)
	}
}
