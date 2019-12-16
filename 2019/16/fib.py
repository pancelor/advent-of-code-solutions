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

N = 10
for n in range(N):
  print n, "|",
  for r in range(n):
    print ncr(n,r),
  print 1

# for k in range(100):
#   print ncr(99+k, 99) % 10,

# for k in range(10000):
#   v = MPowerTopRowIndex(100,k) % 10
#   if k > 20 and v != 0:
#     print k, v
