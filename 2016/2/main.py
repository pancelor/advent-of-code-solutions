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

def clamp(x,a,b):
	return max(a,min(x,b))

def nonedict(d):
	res=defaultdict(lambda: None)
	for k,v in d.items():
		res[k]=v
	return res

labels=nonedict({
	(2,0): '1',
	(1,1): '2',
	(2,1): '3',
	(3,1): '4',
	(0,2): '5',
	(1,2): '6',
	(2,2): '7',
	(3,2): '8',
	(4,2): '9',
	(1,3): 'A',
	(2,3): 'B',
	(3,3): 'C',
	(2,4): 'D',
})

def walk(cmds):
	facing=1
	x=0
	y=2
	# print x,y
	for t in cmds:
		if t=="\n":
			yield labels[(x,y)]
		else:
			px,py=x,y
			dx,dy={
				"U": (0,-1),
				"D": (0,1),
				"L": (-1,0),
				"R": (1,0),
			}[t]
			# print t
			x+=dx #clamp(x+dx,0,4)
			y+=dy #clamp(y+dy,0,4)
			if labels[(x,y)] is None:
				# print 'revert'
				x,y=px,py
			# print x,y

tokens=sys.stdin.read()
# print ''.join(walk(tokens))
for label in walk(tokens):
	print label,
