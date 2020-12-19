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

def parse(text):
	rules={}
	index={}
	p=Parser(text)
	while not p.maybe(r"\n"):
		rulenum,=p.parse(r"(\d+):")
		rulenum=int(rulenum)
		if rulenum not in index:
			index[rulenum]=set()
		opts=[]
		starts=set()
		ns=0
		if p.maybe(r" \""):
			s,=p.parse(r"(\w)\"\n")
			starts.add(s)
			ns=1
		else:
			opts=[]
			while 1:
				nums,_,_=p.parse(r"(( (\d+))+)")
				opts.append(map(int,nums.strip().split(" ")))
				if p.maybe(r"\n"):
					break
				else:
					p.parse(r" \|")
		rules[rulenum]={"opts":opts,"starts":starts,"ns":ns}
		for opt in opts:
			if opt[0] not in index:
				index[opt[0]]=set()
			index[opt[0]].add(rulenum)
	strs=[]
	while not p.done():
		line,=p.parse(r"(\w+)\n")
		strs.append(line)
	return rules,strs,index

def print_rules(rules):
	print "RULES"
	for rid,r in rules.items():
		opts=r["opts"]
		starts=r["starts"]
		print "{}: {} ({})".format(rid,opts,''.join(list(starts)))

rules,strs,index=parse(sys.stdin.read())
pp(index)

# figure out starts
front=set()
seen=set()
for k,r in rules.items():
	if len(r["starts"])>0:
		front.add(k)

while len(front)>0:
	f=front.pop()
	seen.add(f)
	starts=rules[f]["starts"]
	for rid in index[f]:
		r=rules[rid]
		r["starts"].update(starts)
		r["ns"]+=1
		if r["ns"]==len(r["opts"]):
			assert(rid not in seen)
			front.add(rid)
print_rules(rules)

def test_opt(s,i,opt):
	arr1=[i]
	for rid in opt:
		arrtemp=[]
		for itemp in arr1:
			for i2 in test(s,itemp,rid):
				arrtemp.append(i2)
		arr1=arrtemp
	for i2 in arr1:
		yield i2

def test(s,i=0,rid=0):
	r=rules[rid]
	opts=r["opts"]
	starts=r["starts"]
	if i>=len(s):# or s[i] not in starts:
		return
	if len(opts)>0:
		for o in opts:
			for i2 in test_opt(s,i,o):
				yield i2
	else:
		# base case; we already know the first character is in starts
		if s[i]==list(starts)[0]:
			yield i+1

# sys.exit(0)

n=0
for s in strs:
	# print s
	for i in test(s):
		if i==len(s):
			n+=1
			# print "  PASS"
			break
		else:
			pass
			# print "  FAIL (not whole match)"
	# print
print n
