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
	for cc,row in enumerate(grid):
		for rr,entry in enumerate(row):
			print "#" if entry else ".",
		print

grid=[]
W=None
for line in sys.stdin:
	row=[]
	for ch in line.strip():
		row.append(ch=="#")
	grid.append(row)
	if W:
		assert(len(row)==W)
	else:
		W=len(row)

# print_grid(grid)

def p1(dx,dy):
	x=0
	y=0
	trees=0
	while y<len(grid):
		trees+=1 if grid[y][x] else 0
		x+=dx
		y+=dy
		x%=W
	return trees

a=[p1(1,1),p1(3,1),p1(5,1),p1(7,1),p1(1,2)]
print a
print reduce(op.mul,a)
