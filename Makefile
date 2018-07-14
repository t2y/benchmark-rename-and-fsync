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
	go test . -v -testDir $(UNITTEST_DIR) -syncFileRange 102400 -syncFileRangeFlag write

test-bench:
	rm -rf $(BENCHTEST_DIR)
	./main -testDir $(BENCHTEST_DIR) -concurrent 32 -duration 1s -size 1 -benchmark nosync
	./main -testDir $(BENCHTEST_DIR) -concurrent 32 -duration 1s -size 1 -benchmark fsync
	./main -testDir $(BENCHTEST_DIR) -concurrent 32 -duration 1s -size 1 -benchmark fsyn+fadv

DATETIME := $(shell date +"%Y%m%d%H%M%S")

bench:
	sync
	iostat -ymxt 1 /dev/sdk > ${IOSTAT_LOGFILE} &
	# fsync + fadvice
	@date +"%Y%m%d%H%M%S"
	./main -testDir ${DIR_FADV} -concurrent ${CONCURRENT} -duration ${DURATION} -size ${SIZE} -dirMaker ${DIR_MAKER} -benchmark fsyn+fadv
	@date +"%Y%m%d%H%M%S"
	@echo
	sync
	@echo
	# without fsync
	@date +"%Y%m%d%H%M%S"
	./main -testDir ${DIR_NOSYNC} -concurrent ${CONCURRENT} -duration ${DURATION} -size ${SIZE} -dirMaker ${DIR_MAKER} -benchmark nosync
	@date +"%Y%m%d%H%M%S"
	@pkill iostat
	@echo ${IOSTAT_LOGFILE}
	@date +"%Y%m%d%H%M%S"

clean:
	@rm -f main
