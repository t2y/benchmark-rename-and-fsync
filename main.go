package main

import (
	"context"
	"flag"
	"log"
	"math"
	"sync"
	"time"
)

const (
	defaultSyncFileRange = 1048576 // 1MiB
)

var (
	size              int
	testDir           string
	argSyncFileRange  int
	flagSyncFileRange string
	benchmark         string

	concurrent int
	duration   time.Duration
	verbose    bool
)

type benchmarkFunc func(context.Context, int) int

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

	flag.IntVar(&concurrent, "concurrent", 2, "number of goroutines")
	flag.DurationVar(&duration, "duration", 3*time.Second, "run benchmark (e.g. 10s, 1m)")
	flag.BoolVar(&verbose, "verbose", false, "set verbose mode")
}

func main() {
	initFlags()
	flag.Parse()

	benchmarkFunc := getBenchmark()

	resultCh := make(chan int, concurrent)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	startTime := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < int(concurrent); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			numberOfCreated := benchmarkFunc(ctx, i)
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

	log.Printf("total number of created: %d", total)
	log.Printf("- files/second: %v", math.Floor(filesPerSec+0.5))
	log.Printf("- nanoseconds/file: %v", nanoSecPerFile)
}
