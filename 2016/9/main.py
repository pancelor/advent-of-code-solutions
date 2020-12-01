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

def decomp_nextn(src,n):
	# reads the next n bytes of src, and returns the number of bytes it would
	# take up when decompressed

	# read chars
	# if (len x factor) found
	# recurse
	# skip 0 chars (recursing should have handled it properly)

	state=0
	length=0
	factor=0
	res=0
	n0=n
	while n>0:
		ch=next(src)
		assert(ch)
		# print "ch={},n={}/{}".format(ch,n,n0)
		if state==0:
			# 0 normal reading
			n-=1
			if ch=='(':
				state=1
				length=0
				factor=0
			else:
				res+=1
		elif state==1:
			# 1 reading RLE num 1 (len)
			n-=1
			if ch=='x':
				state=2
			else:
				length*=10
				length+=int(ch)
		elif state==2:
			# 2 reading RLE num 2 (num repeats)
			n-=1
			if ch==')':
				# read RLE body
				res+=decomp_nextn(src,length)*factor
				n-=length
				state=0
			else:
				factor*=10
				factor+=int(ch)
		else:
			assert(0)
	# print
	return res


# for line in sys.stdin:
# 	line=line.strip()
# 	print decomp_nextn(iter(line),len(line))

line=sys.stdin.read().strip()
print decomp_nextn(iter(line),len(line))
