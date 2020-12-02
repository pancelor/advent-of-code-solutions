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

def getlines():
	while 1:
		try:
			line=raw_input()
		except EOFError:
			return
		m=re.match(r"^(\d+)\-(\d+) ([a-z]): ([a-z]*)$",line)
		assert(m)
		m0,m1,target,pwd=m.groups()
		yield int(m0),int(m1),target,pwd

def nonedict(d):
	res=defaultdict(lambda: None)
	for k,v in d.items():
		res[k]=v
	return res

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

# acc=0
# for m0,m1,target,pwd in getlines():
# 	c=pwd.count(target)
# 	v=m0 <= c <= m1
# 	if v:
# 		acc+=1
# print acc

acc=0
for m0,m1,target,pwd in getlines():
	c0=pwd[m0-1]
	c1=pwd[m1-1]
	p0=c0==target
	p1=c1==target
	v=p0^p1
	# print c0,c1,p0,p1,v
	if v:
		acc+=1
print acc
