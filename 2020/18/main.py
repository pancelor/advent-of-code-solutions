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

def eval_expr1(p):
	n1=eval_expr2(p)
	if p.maybe(r"\s*\*"):
		n2=eval_expr1(p)
		return n1*n2
	else:
		return n1

def eval_expr2(p):
	n1=eval_expr3(p)
	if p.maybe(r"\s*\+"):
		n2=eval_expr2(p)
		return n1+n2
	else:
		return n1

def eval_expr3(p):
	if p.maybe(r"\s*\("):
		res=eval_expr1(p)
		p.parse(r"\)")
		return res
	else:
		n,=p.parse(r"\s*(-?\d+)")
		return int(n)

num=0
for line in sys.stdin:
	p=Parser(line)
	num+=eval_expr1(p)
print num
