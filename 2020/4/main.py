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

def getline():
	line=raw_input()
	m=re.match(r"^(\d+)$",line)
	assert(m)
	a=m.groups()
	return

def nonedict(d):
	res=defaultdict(lambda: None)
	for k,v in d.items():
		res[k]=v
	return res

def clamp(x,a,b):
	return max(a,min(x,b))

def fieldspresent(passport):
	for key in "byr iyr eyr hgt hcl ecl pid".split(" "):
		if not key in passport:
			# print "invalid; no %s"%key
			return False
	return True

def fieldsvalid(passport):
	if not 1920<=int(passport["byr"])<=2002:
		return False
	if not 2010<=int(passport["iyr"])<=2020:
		return False
	if not 2020<=int(passport["eyr"])<=2030:
		return False
	if not re.match(r"^#[0-9a-f]{6}$", passport["hcl"]):
		return False
	if not re.match(r"^(amb|blu|brn|gry|grn|hzl|oth)$", passport["ecl"]):
		return False
	if not re.match(r"^\d{9}$", passport["pid"]):
		return False
	hgt=re.match(r"^(\d+)(in|cm)$", passport["hgt"])
	if not hgt:
		return False
	units=int(hgt.group(1))
	if hgt.group(2)=="in":
		if not 59<=units<=76:
			return False
	else:
		if not 150<=units<=193:
			return False
	return True

def valid(passport):
	return fieldspresent(passport) and fieldsvalid(passport)

def parse(src):
	passport={}
	for line in src:
		# print("processing",line)
		m=re.match(r"^(|([^:]+):([^:]+))$",line.strip())
		assert(m)
		_,key,value=m.groups()
		if key is None:
			yield passport
			# print "new passport:"
			passport={}
		else:
			if key in passport:
				assert(0)
			passport[key]=value
	yield passport

# I manually modified the input to change spaces to newlines ;)
num=0
for passport in parse(sys.stdin):
	# print passport
	if valid(passport):
		num+=1
print num
