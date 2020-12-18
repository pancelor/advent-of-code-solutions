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

def eval_expr(p):
	n,=p.parse(r"\s*(-?\d+|\()")
	if n=="(":
		acc=eval_expr(p)
		p.parse(r"\)")
	else:
		acc=int(n)
	while not p.done() and not p.peek(r"\s*\)"):
		op,n=p.parse(r"\s*([+*])\s*(-?\d+|\()")
		if n=="(":
			n=eval_expr(p)
			p.parse(r"\)")
		else:
			n=int(n)
		if op=="+":
			acc+=n
		elif op=="*":
			acc*=n
	return acc

num=0
for line in sys.stdin:
	p=Parser(line)
	num+=eval_expr(p)
print num
