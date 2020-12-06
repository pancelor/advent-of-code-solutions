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

def test(passport):
	for key in "byr iyr eyr hgt hcl ecl pid".split(" "):
		if not key in passport:
			# print "invalid; no %s"%key
			return False
	return True

def parse(text):
	p=Parser(text)
	cnt=Counter()
	n=0
	while not p.done():
		n+=1
		letters,=p.parse(r"([a-z]+)\n")
		for c in letters:
			cnt[c]+=1
		print cnt,len(cnt)
		if p.maybe("\n"):
			m=0
			for k,v in cnt.items():
				if v==n:
					m+=1
			print m
			yield m
			n=0
			cnt=Counter()

print sum(parse(sys.stdin.read()))
