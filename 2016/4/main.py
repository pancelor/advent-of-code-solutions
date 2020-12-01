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

def parse(line):
	m=re.match(r"^([a-z\-]+)(\d+)\[([a-z]+)\]$",line)
	assert(m)
	text,num,chk=m.groups()
	num=int(num)
	return text,num,chk

def valid(text,num,chk):
	c=Counter()
	for ch in text:
		if ch!='-':
			c[ch]+=1
	# print list(c.items())
	sk=sorted(c.items(),key=lambda (val,amt): -(amt*1000-ord(val)))
	# print sk
	sk2=''.join(k for k,v in sk[:5])
	# print sk2
	return sk2==chk

# acc=0
# for line in sys.stdin:
# 	text,num,chk=parse(line)
# 	if valid(text,num,chk):
# 		# print text,num,chk
# 		acc+=num
# print acc

def rot(text,num):
	return ''.join(
		' ' if ch=="-"
		else chr(ord('a') + ((ord(ch)-ord('a'))+num)%26) for ch in text)

for line in sys.stdin:
	text,num,chk=parse(line)
	if valid(text,num,chk):
		# print text,num,chk
		print rot(text,num),num
