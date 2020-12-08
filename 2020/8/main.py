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

def sim(prog):
	acc=0
	pc=0
	seen=[False]*len(prog)
	while pc<len(prog):
		if seen[pc]:
			return acc, False
		seen[pc]=True
		cmd=prog[pc]["cmd"]
		num=prog[pc]["num"]
		if cmd=="nop":
			pc+=1
		elif cmd=="acc":
			acc+=num
			pc+=1
		elif cmd=="jmp":
			pc+=num
		else:
			assert(0)
	return acc, True

p=Parser(sys.stdin.read())
prog=[]
while not p.done():
	cmd,num,a=p.parse(r"(nop|acc|jmp) ((-|\+)\d+)\n")
	prog.append({"cmd":cmd,"num":int(num)})

# print prog

for line in prog:
	if line["cmd"]=="acc":
		continue
	line["cmd"]="jmp" if line["cmd"]=="nop" else "nop"
	acc,terminate=sim(prog)
	line["cmd"]="jmp" if line["cmd"]=="nop" else "nop"
	# print line, acc, terminate
	if terminate:
		print acc
