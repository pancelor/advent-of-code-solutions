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

def parse_ticket(p):
	nums=[]
	while 1:
		n,=p.parse(r"(\d+)")
		nums.append(int(n))
		if p.maybe(r","):
			pass
		else:
			p.parse(r"\n")
			break
	return nums

def parse_ranges(p):
	d={}
	while not p.done():
		if p.maybe(r"\n"):
			break
		name,=p.parse(r"([^:]+): ")
		ranges=[]
		while 1:
			a,b=p.parse(r"(\d+)\-(\d+)")
			ranges.append((int(a),int(b)))
			if p.maybe(r"\n"):
				break
			else:
				p.parse(r" or ")
		d[name]=ranges
	return d

def checkany(lookup,val):
	for key in lookup:
		for a,b in lookup[key]:
			if a<=val and val<=b:
				return True
	return False

p=Parser(sys.stdin.read())
lookup=parse_ranges(p)
# print lookup

p.parse(r"your ticket:\n");
myticket=parse_ticket(p)
# print myticket

p.parse(r"\nnearby tickets:\n");
tickets=[]
while not p.done():
	ticket=parse_ticket(p)
	tickets.append(ticket)
# print tickets

n=0
for t in tickets:
	valid=True
	for v in t:
		# print "ticket val",v
		if not checkany(lookup,v):
			valid=False
			n+=v
			break
	print t, valid
print n
