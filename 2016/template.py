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

def walk(cmds):
	facing=1
	x=0
	y=0
	for t in cmds:
		if t=="\n":
			yield x,y
		else:
			dx,dy={
				"U": (0,-1),
				"D": (0,1),
				"L": (-1,0),
				"R": (1,0),
			}[t]
			# print t,dx,dy
			x=clamp(x+dx,-1,1)
			y=clamp(y+dy,-1,1)

tokens=sys.stdin.read()
for x,y in walk(tokens):
	print (y+1)*3+(x+1)+1,
