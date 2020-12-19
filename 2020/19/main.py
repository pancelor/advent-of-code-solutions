#!/usr/bin/env python

import sys
from pprint import pprint as pp
from collections import defaultdict
import re

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
	p=Parser(text)
	rules={}
	index=defaultdict(lambda: set()) # I didn't end up using this
	# parse rules
	while not p.maybe(r"\n"):
		rulenum,=p.parse(r"(\d+):")
		rulenum=int(rulenum)
		opts=[]
		lit=""
		if p.maybe(r" \""):
			s,=p.parse(r"(\w)\"\n")
			lit=s
		else:
			opts=[]
			while 1:
				nums,_,_=p.parse(r"(( (\d+))+)")
				opts.append(map(int,nums.strip().split(" ")))
				if p.maybe(r"\n"):
					break
				else:
					p.parse(r" \|")
		rules[rulenum]={"opts":opts,"lit":lit}
		for opt in opts:
			index[opt[0]].add(rulenum)
	strs=[]
	# parse strings
	while not p.done():
		line,=p.parse(r"([^\n]+)\n")
		strs.append(line)
	return rules,strs,index

def print_rules(rules):
	print "RULES"
	for rid,r in rules.items():
		opts=r["opts"]
		lit=r["lit"]
		print "{}: {}".format(rid,opts or lit)

def match_rule_sequence(rules,s,i,opt):
	"""
	s is a string, opt is an array of rule ids.
	return all indices j such that s[i:j] matches opt (chained matches all in a row)
	"""
	arr1=[i]
	for rid in opt:
		# figure out all possible i2 s.t. s[i:i2] matches all rules up to this point
		arrtemp=[]
		for itemp in arr1:
			for i2 in match_rule(rules,s,itemp,rid):
				arrtemp.append(i2)
		arr1=arrtemp
	for i2 in arr1:
		yield i2

def match_rule(rules,s,i=0,rid=0):
	"""
	s is a string, rid is a rule id.
	return all indices j such that s[i:j] matches the rule with id rid
	"""
	r=rules[rid]
	opts=r["opts"]
	lit=r["lit"]
	if i>=len(s):
		return
	if len(opts)>0:
		for o in opts:
			for i2 in match_rule_sequence(rules,s,i,o):
				yield i2
	else:
		# base case
		if s[i]==lit:
			yield i+1

rules,strs,_=parse(sys.stdin.read())
print_rules(rules)

n=0
for s in strs:
	# print s
	for i in match_rule(rules,s):
		if i==len(s):
			n+=1
			# print "  PASS"
			break
print n
