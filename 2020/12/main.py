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

x=0
y=0
dx=10
dy=-1
dir=0
p=Parser(sys.stdin.read())
while not p.done():
	cmd,num=p.parse(r"(N|S|E|W|L|R|F)(-?\d+)\n")
	num=int(num)
	if cmd=="E":
		dx+=num
	elif cmd=="N":
		dy-=num
	elif cmd=="W":
		dx-=num
	elif cmd=="S":
		dy+=num
	elif cmd=="L" or cmd=="R":
		assert(num==90 or num==180 or num==270)
		if cmd=="L" and num==90 or cmd=="R" and num==270:
			dx,dy=dy,-dx
		elif cmd=="L" and num==180 or cmd=="R" and num==180:
			dx,dy=-dx,-dy
		elif cmd=="L" and num==270 or cmd=="R" and num==90:
			dx,dy=-dy,dx
	elif cmd=="F":
		x+=dx*num
		y+=dy*num
	# print x,y,dx,dy
print x,y
print abs(x)+abs(y)
