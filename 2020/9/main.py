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

nums=[]
for line in sys.stdin:
	nums.append(int(line))

# print nums

n=25
# n=5
target=0
for i in range(len(nums)-n):
	goal=nums[i+n]
	ok=False
	for a,b in itt.combinations(nums[i:i+n],2):
		if a+b==goal:
			ok=True
			break
	if not ok:
		# print i, nums[i:i+n],goal
		target=goal
		break
print target

for i in range(len(nums)):
	for j in range(len(nums)):
		if sum(nums[i:j])==target:
			print min(nums[i:j])+max(nums[i:j])
			sys.exit(0)
