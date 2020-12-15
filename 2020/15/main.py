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

def lasttwo(x):
	found1=False
	for i,y in enumerate(nums)[::-1]:
		if x==y:
			if found1:
				return found1-i
			else:
				found1=i
	return 0

def foo(nums):
	mem=defaultdict(lambda: 0)
	for i in xrange(len(nums)-1):
		mem[nums[i]]=i+1
	last=nums[-1]
	# print "mem",mem.items()
	for i in itt.count(len(nums)):
		# ith number
		previ=mem[last]
		mem[last]=i
		# print "last,i,previ",last,i,previ
		if previ==0:
			last=0
		else:
			last=i-previ
		yield last

# nums=[0,3,6] # the 30000000th number spoken is 175594.
# nums=[1,3,2] # the 30000000th number spoken is 2578.
# nums=[2,1,3] # the 30000000th number spoken is 3544142.
# nums=[1,2,3] # the 30000000th number spoken is 261214.
# nums=[2,3,1] # the 30000000th number spoken is 6895259.
# nums=[3,2,1] # the 30000000th number spoken is 18.
# nums=[3,1,2] # the 30000000th number spoken is 362.

nums=[11,0,1,10,5,19]

i=0
for i,x in enumerate(foo(nums)):
	j=i+len(nums)+1
	if j%3000000==0:
		print "10%"
	if j==30000000:
		print x
		break
