#!/usr/bin/env python

import math
import operator as op
from functools import reduce
import itertools as itt
# https://docs.python.org/2/library/itertools.html
import sys
from pprint import pprint as pp
import re

def walk(cmds):
	acc=0
	facing=1
	x=0
	y=0
	yield x,y
	for t in cmds:
		df=0
		dist=0
		match=re.match(r"^(L|R)(\d+)$",t);
		assert(match)
		if match.group(1)=="L":
			df=1
		else:
			df=-1
		dist=int(match.group(2))
		facing=(facing+df)%4
		for i in range(dist):
			if facing==0:
				x+=1
			elif facing==1:
				y-=1
			elif facing==2:
				x-=1
			elif facing==3:
				y+=1
			yield x,y

# 1
# cmds=sys.stdin.read().split(", ")
# x,y=list(walk(cmds))[-1]
# print abs(x)+abs(y)

# 2
cmds=sys.stdin.read().split(", ")
s=set()
for p in walk(cmds):
	if p in s:
		x,y=p
		print x,y
		print abs(x)+abs(y)
		break
	s.add(p)

# 2 fast:
# some sort of line intersection algorithm
# haven't i done this problem before? on inst laptop maybe?
