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
import md5

def nonedict(d):
	res=defaultdict(lambda: None)
	res.update(d)
	return res

def clamp(x,a,b):
	return max(a,min(x,b))

#
#
#

def hash(n):
	m0 = md5.new()
	m0.update("abc")
	# m0.update("qzyelonm")
	m0.update("%d"%n)
	return m0.hexdigest()

hashes=[]
quints={}
for i in itt.count(0):
	h=hash(i)
	m=re.search(r"(.)\1\1\1\1",h)
	if m:
		quints[i]=m.group(1)
	assert(not re.search(r"(.)\1\1\1\1.*(.)\2\2\2\2",h))
	# print i,h
	hashes.append(h)
	if i%1000==0:
		print i
	if i==100000:
		break
