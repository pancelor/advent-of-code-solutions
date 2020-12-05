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

class Bot:
	def __init__(self,id):
		self.id=id
		self.vals=[]
		self.lo_edge=None
		self.hi_edge=None
	def __str__(self):
		return "#{}{}".format(self.id,self.vals)
	def n(self):
		return len(self.vals)
	def give(self,val):
		self.vals.append(val)
		assert(len(self.vals)<=2)
		self.vals.sort()
		# print "wip {}".format(self)
		if 17 in self.vals and 61 in self.vals:
			print "WINNER {}".format(self.id)
	def takehi(self):
		return self.vals.pop()
	def takelo(self):
		return self.vals.pop(0)
	def set_lo_edge(self,id):
		self.lo_edge=get_bot(id)
	def set_hi_edge(self,id):
		self.hi_edge=get_bot(id)

bots={}
def get_bot(i):
	if i not in bots:
		bots[i]=Bot(i)
	return bots[i]

def parse(text):
	p=Parser(text)
	graph={}
	while not p.done():
		if p.peek(r"value"):
			val,bot_id=p.parse(r"value (\d+) goes to bot (\d+)\n")
			val=int(val)
			bot=get_bot(int(bot_id))
			bot.give(val)
		elif p.peek(r"bot"):
			bot_id,lo_type,lo_id,hi_type,hi_id=p.parse(r"bot (\d+) gives low to (bot|output) (\d+) and high to (bot|output) (\d+)\n")
			bot=get_bot(int(bot_id))
			bot.set_lo_edge(int(lo_id) if lo_type=="bot" else -1-int(lo_id))
			bot.set_hi_edge(int(hi_id) if hi_type=="bot" else -1-int(hi_id))

parse(sys.stdin.read())
fronteir=[]
for b in bots.values():
	# print b
	if b.n()==2:
		fronteir.append(b)

print "ready"

while len(fronteir):
	bot=fronteir.pop()
	lo=bot.lo_edge
	lo.give(bot.takelo())
	if lo.n()==2:
		fronteir.append(lo)
	hi=bot.hi_edge
	hi.give(bot.takehi())
	if hi.n()==2:
		fronteir.append(hi)

# for b in bots.values():
# 	print b
print get_bot(-1).takelo()*get_bot(-2).takelo()*get_bot(-3).takelo()
