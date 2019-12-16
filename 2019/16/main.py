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

vals = map(int, raw_input())
offset = listToDec(vals[:7])
# print offset

# print len(vals)
res = stepMany(vals, 100)
print listToDec(res[:8])
# print res
