#!/usr/bin/env python3

import math
import sys

def get_percentile(vals, pct):
	if len(vals) == 0:
		return float('nan')
	k = float(len(vals)-1) * pct
	f = math.floor(k)
	c = math.ceil(k)
	if f == c:
		return vals[int(k)]
	d0 = vals[int(f)] * (c - k)
	d1 = vals[int(c)] * (k - f)
	return d0 + d1

all = []
min = float('inf')
max = float('-inf')
count = 0
total = float(0)
for line in sys.stdin:
	count += 1
	v = float(line)
	total += v
	all.append(v)
	if v > max:
		max = v
	if v < min:
		min = v

all.sort()
pct25 = get_percentile(all, 0.25)
median = get_percentile(all, 0.5)
pct75 = get_percentile(all, 0.75)
mean = total / count

stddev = sum((x-mean)**2 for x in all)
stddev /= len(all)-1
stddev = stddev**0.5

print("sum:    " + str(total))
print("min:    " + str(min))
print("pct25:  " + str(pct25))
print("median: " + str(median))
print("pct75:  " + str(pct75))
print("max:    " + str(max))
print("mean:   " + str(mean))
print("stddev: " + str(stddev))
print("count:  " + str(count))
