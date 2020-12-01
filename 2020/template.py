#!/usr/bin/env python

import operator as op
from functools import reduce
import itertools as itt
# https://docs.python.org/2/library/itertools.html
import sys
from pprint import pprint as pp

def p1(nums,goal):
	d={}
	for val in nums:
		pair=goal-val
		if pair in d:
			return [pair,val]
		d[val]=pair

def p2(nums,goal):
	s=set(nums)
	for a,b in itt.combinations(nums,2):
		c=goal-a-b
		if c in s:
			return a,b,c

nums=map(int,sys.stdin)

p,v=p1(nums,2020)
print "part 1:",p,v
print p*v

# p,v,n=p2(nums,2020)
# print "part 2:",p,v,n
# print p*v*n
