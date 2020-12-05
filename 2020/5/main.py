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

# I preprocessed the input in my text editor
# for this problem: F->0 B->1 L->0 R->1

ids=[int(line.strip(),2) for line in sys.stdin]
# print max(ids)
print sorted(ids)
# i=54
# for x in sorted(ids):
# 	if x!=i:
# 		print i
# 		break
# 	i+=1
