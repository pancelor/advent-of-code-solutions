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
		self._last=None
	def done(self):
		"""
		returns true iff the only stuff left is spaces
		"""
		return self.peek(r"\s*$")
	def peek(self,rgx):
		"""
		returns the match, without mutating internal state
		"""
		self._last=re.match(rgx,self.text)
		return self._last
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
		self._last=m
		return m
	def last(self):
		"""
		returns the last thing matched
		"""
		return self._last
	def parse(self,rgx):
		"""
		returns the _groups_ and mutates state.
		_requires_ a match to be found
		"""
		m=self.maybe(rgx)
		assert(m)
		self._last=m
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

registers={"a":0,"b":0,"c":0,"d":0}
def getvalue(str):
	if re.match(r"-?\d+$",str):
		return int(str)
	elif re.match(r"[a-d]$",str):
		return registers[str]

def print_registers():
	print "%d|%d|%d|%d"%(
		registers["a"],
		registers["b"],
		registers["c"],
		registers["d"])

I_NONE=0
I_CPY=1
I_INC=2
I_DEC=3
I_JNZ=4
def parse(text):
	p=Parser(text)
	while not p.done():
		cmd={"t":I_NONE,"x":None,"y":None}
		if p.peek(r"cpy "):
			x,y=p.parse(r"cpy (-?\d+|[a-d]) (-\d+|[a-d])\n")
			cmd["t"]=I_CPY
			cmd["x"]=x
			cmd["y"]=y
		elif p.peek(r"(inc|dec) "):
			i,x=p.parse(r"(inc|dec) ([a-d])\n")
			cmd["t"]=I_INC if i=="inc" else I_DEC
			cmd["x"]=x
		elif p.peek(r"jnz "):
			x,y=p.parse(r"jnz (-?\d+|[a-d]) (-?\d+)\n")
			cmd["t"]=I_JNZ
			cmd["x"]=x
			cmd["y"]=y
		else:
			assert(0)
		yield cmd

prog=list(parse(sys.stdin.read()))
pc=0
while pc<len(prog):
	cmd=prog[pc]
	# print pc,cmd,
	# print_registers()
	if cmd["t"]==I_INC:
		registers[cmd["x"]]+=1
		pc+=1
	elif cmd["t"]==I_DEC:
		registers[cmd["x"]]-=1
		pc+=1
	elif cmd["t"]==I_CPY:
		registers[cmd["y"]]=getvalue(cmd["x"])
		pc+=1
	elif cmd["t"]==I_JNZ:
		if getvalue(cmd["x"])!=0:
			pc+=int(cmd["y"])
		else:
			pc+=1
	else:
		assert(0)
print_registers()
