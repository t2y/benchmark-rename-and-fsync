#!/bin/bash

TOP_DIR=${1:-tdir}
TMP_DIR="/tmp"
DURATION="2m"

CONCURRENTS=(32 64 128 256 512 1024)
LOCATIONS=("sequential" "random")
SIZES=(1 4 10 100)

echo "top test directory: $TOP_DIR"
if [ ! -d "$TOP_DIR" ]; then
    mkdir -p "$TOP_DIR"
fi

start=$(date +"%Y-%m-%d %H:%M:%S")
echo "start benchmark: $start"

for concurrent in "${CONCURRENTS[@]}"; do
    for loc in "${LOCATIONS[@]}"; do
        for size in "${SIZES[@]}"; do
            prefix="${loc}-con${concurrent}-${size}KB"
            echo "running $prefix ..."
            now=$(date +"%Y%m%d%H%M%S")
            iostat_logfile="${TMP_DIR}/iostat-${prefix}-${now}.log"
            dir_fadv="${TOP_DIR}/fadv-${prefix}"
            dir_nosync="${TOP_DIR}/nosync-${prefix}"
            make IOSTAT_LOGFILE="$iostat_logfile" DIR_FADV="$dir_fadv" DIR_NOSYNC="$dir_nosync" CONCURRENT=$concurrent DURATION=$DURATION SIZE=$size DIR_MAKER=$loc bench;
            echo
            echo
        done
    done
done

end=$(date +"%Y-%m-%d %H:%M:%S")
echo "end benchmark: $end"
