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

def nonedict(d):
	res=defaultdict(lambda: None)
	res.update(d)
	return res

def clamp(x,a,b):
	return max(a,min(x,b))

def get(grid,x,y):
	if 0<=x and x < len(grid[0]) and 0<=y and y<len(grid):
		return grid[y][x]

def count(grid,x,y):
	num=0
	for dx in [-1,0,1]:
		for dy in [-1,0,1]:
			if dx==0 and dy==0: continue
			i=0
			while True:
				i+=1
				c=get(grid,x+i*dx,y+i*dy)
				if c is None:
					break
				if c=="L":
					break
				if c=="#":
					num+=1
					break

	return num

def frob(grid):
	changes=False
	grid2=[]
	for y,row in enumerate(grid):
		row2=[v for v in row]
		for x,val in enumerate(row):
			nocc = count(grid,x,y)
			if val=="L" and nocc==0:
				row2[x]="#"
				changes=True
			elif val=="#" and nocc>=5:
				row2[x]="L"
				changes=True
		grid2.append(row2)

	return grid2,changes

grid=[]
for line in sys.stdin:
	grid.append(list(line.strip()))
# pp(grid)

for i in itt.count():
	grid,changes=frob(grid)
	# pp(list(''.join(line) for line in grid))
	if not changes:
		print i
		break

n=0
for row in grid:
	for val in row:
		if val=="#":
			n+=1
print n
