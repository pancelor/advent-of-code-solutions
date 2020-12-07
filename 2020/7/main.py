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
	rule={}
	while not p.done():
		first,=p.parse(r"(\w+ \w+) bags contain ")
		contents=[]
		while 1:
			if p.maybe(r"no other bags\.\n"):
				break
			contents.append(p.parse(r"(\d+) (\w+ \w+) bags?"))
			if p.maybe(r"\.\n"):
				break
			else:
				p.parse(r", ")
		rule["first"]=first
		rule["contents"]=contents
		yield rule
		rule={}

rules=list(parse(sys.stdin.read()))

# allowed=set(["shiny gold"])
# seen=set()
# num=0
# while len(allowed)>0:
# 	t=allowed.pop()
# 	seen.add(t)
# 	num+=1
# 	print t
# 	for r in rules:
# 		for c in r["contents"]:
# 			if c[1]==t and r["first"] not in seen:
# 				allowed.add(r["first"])
# print num-1

def thing(name):
	n=1
	for r in rules:
		if r["first"]==name:
			for c in r["contents"]:
				m,other=c
				n+=int(m)*thing(other)
			print r,n
	return n

print thing("shiny gold")-1
