#!/usr/bin/env python

import operator as op
from functools import reduce
import itertools as itt
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
	for n in nums:
		x=p1(nums,goal-n)
		if x:
			p,v=x
			return p,v,n

def p2fast(nums,goal):
	s=set(nums)
	for a,b in itt.combinations(nums,2):
		c=goal-a-b
		if c in s:
			return a,b,c

nums=map(int,sys.stdin)

# p,v=p1(nums,2020)
# print "part 1:",p,v
# print p*v

# p,v,n=p2(nums,2020)
# print "part 2 (optimized):",p,v,n
# print p*v*n

p,v,n=p2fast(nums,2020)
print "part 2 (fast):",p,v,n
print p*v*n
