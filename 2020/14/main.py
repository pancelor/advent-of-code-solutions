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

mem=defaultdict(lambda: 0)
passmask=0
stamp=0
p=Parser(sys.stdin.read())
while not p.done():
	if p.peek(r"mask"):
		mask,_=p.parse(r"mask = ((X|0|1)+)\n")
		passmask=0
		stamp=0
		for i,ch in enumerate(reversed(mask)):
			if ch=="X":
				passmask|=1<<i
			elif ch=="1":
				stamp|=1<<i
		print passmask,stamp
		print("passmask {0:>08b}".format(passmask))
		print("stamp {0:>08b}".format(stamp))
	else:
		addr,val=p.parse(r"mem\[(\d+)\] = (\d+)\n")
		addr=int(addr)
		val=int(val)
		mem[addr]=(val&passmask)|stamp
		print("mem[{0}]={1:>08b}={1}".format(addr,mem[addr]))

sum=0
for k,v in mem.items():
	sum+=v
print sum
