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

def decomp(line):
	state=0
	# 0 normal reading
	# 1 reading RLE num 1 (len)
	# 2 reading RLE num 2 (num repeats)
	# 3 reading RLE body
	resbuf=[]
	num1=0
	num2=0
	bodybuf=[]
	for ch in line:
		if state==0:
			if ch=='(':
				state=1
				num1=0
				num2=0
			else:
				resbuf.append(ch)
		elif state==1:
			if ch=='x':
				state=2
			else:
				num1*=10
				num1+=int(ch)
		elif state==2:
			if ch==')':
				state=3
				bodybuf=[]
			else:
				num2*=10
				num2+=int(ch)
		elif state==3:
			bodybuf.append(ch)
			num1-=1
			if num1==0:
				resbuf+=bodybuf*num2
				state=0
		else:
			assert(0)
	return ''.join(resbuf)

# for line in sys.stdin:
# 	print decomp(line.strip())

arg=sys.stdin.read().strip()
print len(arg)
arg=decomp(arg)
print len(arg)

