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

def validtls(line):
	in_bracket=False
	found=False
	p3=None
	p2=None
	p1=None
	for ch in line:
		if p3==ch and p2==p1 and p2!=p3:
			if in_bracket:
				return False
			else:
				found=True
		if ch=="[":
			assert(not in_bracket)
			in_bracket=True
		if ch=="]":
			assert(in_bracket)
			in_bracket=False
		p3,p2,p1=p2,p1,ch
	return found

def validssl(line):
	in_bracket=False
	partition=[[],[]]
	found=False
	p2=None
	p1=None
	for ch in line:
		if p2==ch and p2!=p1:
			ls=partition[1 if in_bracket else 0].append("%s%s%s"%(p2,p1,ch))
		if ch=="[":
			assert(not in_bracket)
			in_bracket=True
		if ch=="]":
			assert(in_bracket)
			in_bracket=False
		p2,p1=p1,ch
	b=set(partition[1])
	for s in partition[0]:
		s2=s[1]+s[0]+s[1]
		if s2 in b:
			return True
	return False

acc=0
for line in sys.stdin:
	v=validssl(line)
	# print line,v
	if v:
		acc+=1
print acc
