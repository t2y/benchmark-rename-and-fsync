import sys
from datetime import datetime

import matplotlib
matplotlib.use('Agg')
from matplotlib import pyplot as plt


class Plotter:

    def __init__(self, since=None, until=None):
        figsize = (12, 6)
        self.fig = plt.figure(figsize=figsize)
        self.fig.suptitle('profile plotter')
        self.subplot = self.fig.add_subplot(1, 1, 1)
        self.subplot.set_ylabel('milliseconds')

        if since is not None and until is not None:
            self.subplot.set_xlim(since, until)

    def plot(self, dates, elapsed_times):
        self.subplot.plot(dates, elapsed_times, 'k.')

    def save(self):
        plt.savefig('profile_poiner.png')


def parse_line(line):
    # line format: 2018-06-26T00:32:21.500767889+09:00,5.555149\n
    row = line.rstrip().split(',')
    date_str = row[0].split('+')[0][:26]  # strip nanoseconds
    millisec = float(row[1])
    date = datetime.strptime(date_str, '%Y-%m-%dT%H:%M:%S.%f')
    return date, millisec


def main():
    dates, elapsed_times = [], []
    with open(sys.argv[1]) as f:
        for line in f:
            date, elapsed_time = parse_line(line)
            dates.append(date)
            elapsed_times.append(elapsed_time)
    print('number of data', len(dates))

    since, until = None, None
    if len(sys.argv) > 2:
        since = datetime.strptime(sys.argv[2], '%Y%m%d%H%M%S')
        until = datetime.strptime(sys.argv[3], '%Y%m%d%H%M%S')

    plotter = Plotter(since, until)
    plotter.plot(dates, elapsed_times)
    plotter.save()


def usage():
    print('usage: python %s profile.csv' % __file__)
    print('       e.g.) python %s profile.csv' % __file__)


if __name__ == '__main__':
    if len(sys.argv) < 2:
        usage()
        sys.exit(1)
    main()
