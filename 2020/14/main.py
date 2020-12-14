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

class Parser:
	def __init__(self,text):
		self.text=text
	def done(self):
		"""
		returns true iff the only stuff left is spaces
		"""
		return self.peek(r"\s*$")
	def peek(self,rgx):
		"""
		returns the match, without mutating internal state
		"""
		return re.match(rgx,self.text)
	def maybe(self,rgx):
		"""
		returns the match and mutates state if a match was found
		if no match found, nothing changes
		"""
		# match only matches at beginning of text
		m=self.peek(rgx)
		if m:
			a,b=m.span()
			assert(a==0)
			self.text=self.text[b:]
		return m
	def parse(self,rgx):
		"""
		returns the _groups_ and mutates state.
		_requires_ a match to be found
		"""
		m=self.maybe(rgx)
		assert(m)
		return m.groups()

#
#
#

def spread(xmask_bits):
	for i in range(2**len(xmask_bits)):
		n=0
		for j,b_ix in enumerate(xmask_bits):
			if i&(1<<j)>0:
				n|=1<<b_ix
		yield n

mem=defaultdict(lambda: 0)
xmask=0
xmask_bits=[]
ormask=0
p=Parser(sys.stdin.read())
while not p.done():
	if p.peek(r"mask"):
		mask,_=p.parse(r"mask = ((X|0|1)+)\n")
		xmask=0
		xmask_bits=[]
		ormask=0
		for i,ch in enumerate(reversed(mask)):
			if ch=="X":
				xmask|=1<<i
				xmask_bits.append(i)
			elif ch=="1":
				ormask|=1<<i
		# print "ormask {:>08b}".format(ormask)
		# print "xmask {:>08b}".format(xmask)
		# for s in spread(xmask_bits):
		# 	print "  spread {:>08b}".format(s)
	else:
		addr,val=p.parse(r"mem\[(\d+)\] = (\d+)\n")
		addr=int(addr)
		val=int(val)
		combined=(addr|ormask)&(~xmask)
		# print addr,"combined {0:>08b}".format(combined)
		for s in spread(xmask_bits):
			c2=combined|s
			# print "c2 {0:>08b}".format(c2)
			mem[c2]=val

sum=0
for k,v in mem.items():
	sum+=v
print sum
