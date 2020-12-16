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

def checkfield(lookup,key,val):
	for a,b in lookup[key]:
		if a<=val and val<=b:
			return True
	return False

def checkany(lookup,val):
	for key in lookup:
		if checkfield(lookup,key,val):
			return True
	return False

def valid(lookup,t):
	valid=True
	for v in t:
		if not checkany(lookup,v):
			return False
	return True

p=Parser(sys.stdin.read())
lookup=parse_ranges(p)
nfields=len(lookup)
test=nfields<10
if test:
	print lookup

p.parse(r"your ticket:\n");
myticket=parse_ticket(p)
if test:
	print "my", myticket

p.parse(r"\nnearby tickets:\n");
tickets=[]
while not p.done():
	ticket=parse_ticket(p)
	if valid(lookup,ticket):
		tickets.append(ticket)
if test:
	print "nearby",tickets

# poss[i] is the list of possibilities the ith field could be
poss=defaultdict(lambda: lookup.keys())

changes=True
while changes:
	changes=False
	# print "\n\n\nNEW_ITERATION\n\n"
	for t in tickets:
		# print "ticket:",t
		for i,v in enumerate(t):
			for key in poss[i][:]:
				if not checkfield(lookup,key,v):
					# print "throwing out key '%s' for field %d"%(key,i)
					poss[i].remove(key)
					changes=True
			if len(poss[i])<=1:
				assert(len(poss[i])==1)
				key=poss[i][0]
				# print "only possibility left is",key
				for j in range(nfields):
					if i!=j:
						# print "throwing out key '%s' for field %d"%(key,j)
						if key in poss[j]:
							poss[j].remove(key)
							changes=True
# pp(poss)
prod=1
for i,keys in poss.items():
	key=keys[0]
	if re.match(r"departure ",key):
		print i,key
		prod*=myticket[i]
print prod
