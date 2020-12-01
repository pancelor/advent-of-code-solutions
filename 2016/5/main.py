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

p=['_']*8
m0 = md5.new()
m0.update("reyedfim")
# m0.update("abc")
# for i in itt.count(3000000):
for i in itt.count(0):
	if i%500000==0:
		print "i=%d"%i
	m1=m0.copy()
	m1.update("%d"%i)
	res=m1.hexdigest()
	if res[:5]=="00000":
		print res
		pos=res[5]
		posi=ord(pos)-ord('0')
		if 0 <= posi and posi <= 7:
			ch=res[6]
			if p[posi]=='_':
				p[posi]=ch
				print ''.join(p)
				if '_' not in p:
					break
