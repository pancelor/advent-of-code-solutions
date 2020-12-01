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

def getline():
	line=raw_input()
	match=re.match(r"^\s*(\d+)\s+(\d+)\s+(\d+)\s*$",line)
	assert(match)
	a=int(match.group(1))
	b=int(match.group(2))
	c=int(match.group(3))
	return a,b,c

def munge():
	while 1:
		a=""
		try:
			a=getline()
		except EOFError:
			return
		b=getline()
		c=getline()
		yield a[0],b[0],c[0]
		yield a[1],b[1],c[1]
		yield a[2],b[2],c[2]

def valid3(a,b,c):
	return a+b>c and a+c>b and b+c>a

count=0
for a,b,c in munge():
	v=valid3(a,b,c)
	# print a,b,c,v
	if v:
		count+=1
print count
