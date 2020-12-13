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

def foo(p):
	est,=p.parse(r"(\d+)\n")

	est=int(est)
	rem=-1
	while not p.done():
		rem+=1
		n,_=p.parse(r"(\d+|x)(,|\n)")
		# print n
		if n=="x":
			continue
		else:
			n=int(n)
			# wait=-(est%n)%n
			yield [n,(-rem)%n]

p=Parser(sys.stdin.read())
ls=sorted(foo(p),key=lambda x:x[0])

print ls

x=ls[0][1]
m=ls[0][0]
for ix in range(len(ls)-1):
	for i in itt.count():
		# print i,x,m
		if x%ls[ix+1][0]==ls[ix+1][1]:
			print "found!",x,ls[:ix+2]
			m*=ls[ix+1][0]
			break
		x+=m

