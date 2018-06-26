import sys
from collections import defaultdict
from datetime import datetime, timedelta


def calc_throughput(profile):
    date_ranges = sorted(profile.keys())
    prev_date = date_ranges[0]
    end_date = date_ranges[-1]

    window_size = 10
    throughputs = []
    while prev_date < end_date:
        total = 0
        for sec in range(1, 1 + window_size):
            target_date = prev_date + timedelta(seconds=sec)
            total += profile[target_date]  # KB
        tp = round(total / (1024.0 * window_size), 2)
        throughputs.append(tp)
        # print('datetime:', target_date, 'throughput:', tp)
        prev_date = target_date

    print(','.join(str(i) for i in throughputs))


def parse_line(line):
    # line format: 2018-06-26T00:32:21.500767889+09:00,5.555149\n
    row = line.rstrip().split(',')
    date_str = row[0].split('+')[0][:26]  # strip nanoseconds
    millisec = float(row[1])
    start = datetime.strptime(date_str, '%Y-%m-%dT%H:%M:%S.%f')
    end = start + timedelta(milliseconds=millisec)
    end = end.replace(microsecond=0)  # strip microseconds to calculate simply
    return end


def main():
    size = int(sys.argv[2])
    profile = defaultdict(int)
    with open(sys.argv[1]) as f:
        for line in f:
            end_date = parse_line(line)
            profile[end_date] += size

    calc_throughput(profile)


def usage():
    print('usage: python %s path/to/profile.csv file_size' % __file__)
    print('       e.g.) python %s profile.csv 1' % __file__)


if __name__ == '__main__':
    if len(sys.argv) != 3:
        usage()
        sys.exit(1)
    main()
