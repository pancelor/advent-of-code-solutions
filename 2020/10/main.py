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

def getline():
	line=raw_input()
	m=re.match(r"^(\d+)$",line)
	assert(m)
	a=m.groups()
	return

ratings=sorted(list(map(int,sys.stdin)))
N=len(ratings)

# print ratings

c=Counter()
for i in range(N-1):
	a=ratings[i]
	b=ratings[i+1]
	c[b-a]+=1
c[1]+=1
c[3]+=1

# print c, c[1]*c[3]

M=ratings[-1]+3
table=[0]*M
table[0]=1
table[1]=1
for r in ratings:
	# print i
	ways=0
	if r>=1:
		ways+=table[r-1]
	if r>=2:
		ways+=table[r-2]
	if r>=3:
		ways+=table[r-3]
	table[r]=ways
	# print r, table
print table[-3]
