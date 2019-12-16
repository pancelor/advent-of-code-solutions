#!/usr/bin/env python

import itertools as itt

def listToDec(vals):
  total = 0
  for x in vals:
    total *= 10
    total += x
  return total

#
# PART 1
#

# def coeffs(n):
#   res = itt.repeat(0, n)
#   res = itt.chain(res, itt.repeat(1, n))
#   res = itt.chain(res, itt.repeat(0, n))
#   res = itt.chain(res, itt.repeat(-1, n))
#   res = itt.cycle(res)
#   res.next()
#   return res

# def step(vals):
#   for i in range(len(vals)):
#     k = 0
#     for a, b in zip(vals, coeffs(i+1)):
#       k += a*b
#     yield abs(k) % 10

# def stepMany(vals, n):
#   for i in range(n):
#     # print "step", i
#     vals = list(step(vals))
#   return vals

# vals = map(int, raw_input())
# res = stepMany(vals, 100)
# print listToDec(res[:8])

#
# PART 2
#

import operator as op
from functools import reduce

def ncr(n, r):
  r = min(r, n-r)
  numer = reduce(op.mul, range(n, n-r, -1), 1)
  denom = reduce(op.mul, range(1, r+1), 1)
  return numer / denom

def MPowerTopRowIndex(p, i):
  # returns the ith member (0-based) of the top row of M^p,
  # where M is an upper triangular all ones matrix
  # (each digit of the final answer involves a dot product between M^100 and vals)
  return ncr(p-1+i, p-1)

def findDigit(digit_ix, vals, offset, mp100tri_cache):
  total = 0
  for i in itt.count(0):
    supervalsIndex = offset + digit_ix + i
    if supervalsIndex == len(vals) * 10000:
      break
    valsDigit = vals[supervalsIndex % len(vals)]
    total += valsDigit * mp100tri_cache[i]
    total %= 10
  # print "digit #{}: {}".format(digit_ix, total)
  return total

vals = map(int, raw_input())
offset = listToDec(vals[:7])
# print "offset\n", offset

mp100tri_cache = []
N = len(vals) * 10000 - offset+10
for i in range(N):
  if i%100000 == 0:
    print "precomputing ncr {}/{}".format(i, N)
  mp100tri_cache.append(MPowerTopRowIndex(100, i))

print listToDec([findDigit(d, vals, offset, mp100tri_cache) for d in range(8)])
