import itertools as itt

def coeffs(n):
  res = itt.repeat(0, n)
  res = itt.chain(res, itt.repeat(1, n))
  res = itt.chain(res, itt.repeat(0, n))
  res = itt.chain(res, itt.repeat(-1, n))
  res = itt.cycle(res)
  res.next()
  return res

def step(vals):
  for i in range(len(vals)):
    k = 0
    for a, b in zip(vals, coeffs(i+1)):
      k += a*b
    yield abs(k) % 10

def stepMany(vals, n):
  for i in range(n):
    # print "step", i
    vals = list(step(vals))
  return vals

def listToDec(vals):
  total = 0
  for x in vals:
    total *= 10
    total += x
  return total

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
  return ncr(p-1+i, p-1)



vals = map(int, raw_input())
offset = listToDec(vals[:7])
print "offset\n", offset

# first digit
final = 0
for digit in range(8):
  total = 0
  for i in itt.count(0):
    supervalsIndex = offset + digit + i
    if supervalsIndex%100000 == 0:
      print supervalsIndex, "/", len(vals) * 10000
    if supervalsIndex == len(vals) * 10000:
      break
    valsDigit = vals[supervalsIndex % len(vals)]
    total += valsDigit * MPowerTopRowIndex(100, i)
    total %= 10
  print digit, ":", total
  final *= 10
  final += total
print final

# res = stepMany(vals, 100)
# print listToDec(res[:8])
# print res
