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


# def MPowerTopRowIndex(p, i):
#   # returns the ith member (0-based) of the top row of M^p,
#   # where M is an upper triangular all ones matrix
#   # (each digit of the final answer involves a dot product between M^100 and vals)
#   return ncr(p-1+i, p-1)

def findDigit(digit_ix, vals, offset, MPowerTopRowIndex):
  total = 0
  for i in itt.count(0):
    supervalsIndex = offset + digit_ix + i
    if supervalsIndex % 100000 == 0:
      print supervalsIndex
    if supervalsIndex == len(vals) * 10000:
      break
    valsDigit = vals[supervalsIndex % len(vals)]
    total += valsDigit * MPowerTopRowIndex[i]
    total %= 10
  # print "digit #{}: {}".format(digit_ix, total)
  return total

vals = map(int, raw_input())
offset = listToDec(vals[:7])
# print "offset\n", offset

MPowerTopRowIndex = [1]
for i in xrange(1, len(vals) * 10000 - offset):
  MPowerTopRowIndex.append((MPowerTopRowIndex[-1]*(99+i))/i)

print listToDec([findDigit(d, vals, offset, MPowerTopRowIndex) for d in range(8)])
