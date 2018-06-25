import sys
from collections import defaultdict
from datetime import datetime, timedelta


def calc_throughput(profile):
    date_ranges = sorted(profile.keys())
    start_date = date_ranges[0]

    total = 0
    throughputs = []
    for key in date_ranges:
        total += profile[key]  # KB
        t = key - start_date
        if t.seconds % 10 == 0 and t.seconds != 0:
            tp = round(total / (1024.0 * t.seconds), 2)
            throughputs.append(tp)
            print('datetime:', key, 'throughput:', tp)

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
    profile = defaultdict(int)
    with open(sys.argv[1]) as f:
        for line in f:
            end_date = parse_line(line)
            profile[end_date] += 1

    calc_throughput(profile)


def usage():
    print('usage: python %s profile.csv' % __file__)


if __name__ == '__main__':
    if len(sys.argv) != 2:
        usage()
        sys.exit(1)
    main()
