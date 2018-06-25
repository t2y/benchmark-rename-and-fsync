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

NOSYNC_DIR := nosync-testdata
FSYNC_DIR := fsync-testdata
FADV_DIR := fsyn+fadv-testdata

DATETIME := $(shell date +"%Y%m%d%H%M%S")
IOSTAT_LOG := iostat-$(DATETIME).log
SLEEP_TIME := 60

bench: clean-data
	iostat -ymxt 1 /dev/sdk > $(IOSTAT_LOG) &
	# fsync + fadvice
	@date +"%Y%m%d%H%M%S"
	./main -testDir $(FADV_DIR) -concurrent ${CONCURRENT} -duration ${DURATION} -size ${SIZE} -dirMaker ${DIR_MAKER} -benchmark fsyn+fadv
	@date +"%Y%m%d%H%M%S"
	@echo
	sleep $(SLEEP_TIME)
	@echo
	# fsync
	@date +"%Y%m%d%H%M%S"
	./main -testDir $(FSYNC_DIR) -concurrent $(CONCURRENT) -duration $(DURATION) -size ${SIZE} -dirMaker ${DIR_MAKER} -benchmark fsync
	@date +"%Y%m%d%H%M%S"
	@echo
	sleep $(SLEEP_TIME)
	@echo
	# without fsync
	@date +"%Y%m%d%H%M%S"
	./main -testDir $(NOSYNC_DIR) -concurrent $(CONCURRENT) -duration $(DURATION) -size ${SIZE} -dirMaker ${DIR_MAKER} -benchmark nosync
	@date +"%Y%m%d%H%M%S"
	sleep $(SLEEP_TIME)
	@pkill iostat
	@echo $(IOSTAT_LOG)

clean-data:
	@rm -rf $(UNITTEST_DIR)
	@rm -rf $(BENCHTEST_DIR)
	@rm -rf $(NOSYNC_DIR)
	@rm -rf $(FSYNC_DIR)
	@rm -rf $(FADV_DIR)

clean:
	@rm -f main
