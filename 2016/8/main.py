#!/usr/bin/env python

import math
import operator as op
from functools import reduce
import itertools as itt
# https://docs.python.org/2/library/itertools.html
from collections import Counter,defaultdict
import sys
from pprint import pprint as pp
import re

def getline():
	line=raw_input()
	m=re.match(r"^(\d+)$",line)
	assert(m)
	a=m.groups()
	return

def nonedict(d):
	res=defaultdict(lambda: None)
	for k,v in d.items():
		res[k]=v
	return res

def clamp(x,a,b):
	return max(a,min(x,b))

def print_grid(grid):
	for row in grid:
		print ''.join('#' if v else "." for v in row)

W=50
H=6
grid=[]
for y in range(H):
	grid.append([0]*W)

print_grid(grid)
for line in sys.stdin:
	m=re.match(r"^rect (\d+)x(\d+)|rotate row y=(\d+) by (\d+)|rotate column x=(\d+) by (\d+)$",line)
	assert(m)
	if m.group(1):
		w=int(m.group(1))
		h=int(m.group(2))
		print ">\tRECT",w,h
		for y in range(h):
			for x in range(w):
				grid[y][x]=1
	elif m.group(3):
		y=int(m.group(3))
		d=int(m.group(4))
		print ">\tROW",y,d
		newrow=[0]*W
		for x in range(W):
			newrow[x]=grid[y][(x-d)%W]
		for x in range(W):
			grid[y][x]=newrow[x]
	elif m.group(5):
		x=int(m.group(5))
		d=int(m.group(6))
		print ">\tCOL",x,d
		newcol=[0]*W
		for y in range(H):
			newcol[y]=grid[(y-d)%H][x]
		for y in range(H):
			grid[y][x]=newcol[y]
	print_grid(grid)

acc=0
for row in grid:
	acc+=sum(row)
print acc
