#!/bin/bash

TOP_DIR=${1:-tdir}
TMP_DIR="/tmp"
DURATION="2m"
RESULTS_PREFIX="results"

RATES=($(seq 150 50 500))
HANDLERS=("fsync" "nosync")

echo "top test directory: $TOP_DIR"
if [ ! -d "$TOP_DIR" ]; then
    mkdir -p "$TOP_DIR"
fi

start=$(date +"%Y-%m-%d %H:%M:%S")
echo "start requesting: $start"

for rate in "${RATES[@]}"; do
    for handler in "${HANDLERS[@]}"; do
        sync
        echo "request rate is $rate, target is $handler"
        now=$(date +"%Y%m%d%H%M%S")
        iostat_logfile="${TMP_DIR}/iostat-vegeta-${rate}-${handler}-${now}.log"
	    iostat -ymxt 1 /dev/sdk > ${iostat_logfile} &
        echo "GET http://localhost:8090/$handler" | vegeta attack -duration=$DURATION -rate $rate | tee "${RESULTS_PREFIX}-${rate}-${handler}.bin" | vegeta report
	    pkill iostat
        echo
    done
done

end=$(date +"%Y-%m-%d %H:%M:%S")
echo "end requesting: $end"
