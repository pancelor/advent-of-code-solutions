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

counts=[Counter(),Counter(),Counter(),Counter(),Counter(),Counter(),Counter(),Counter(),]
for line in sys.stdin:
	for i,ch in enumerate(line.strip()):
		assert(i<8)
		counts[i][ch]+=1
r=[]
for i in range(8):
	x=counts[i].most_common()[-1:][0][0]
	r.append(x)
print ''.join(r)
