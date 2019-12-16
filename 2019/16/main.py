import itertools as itt

def coeffs(n):
  res = itt.repeat(0, n)
  res = itt.chain(res, itt.repeat(1, n))
  res = itt.chain(res, itt.repeat(0, n))
  res = itt.chain(res, itt.repeat(-1, n))
  res = itt.cycle(res)
  res.next()
  return res

vals = map(int, raw_input())

def step(vals):
  for i in range(len(vals)):
    k = 0
    for a, b in zip(vals, coeffs(i+1)):
      k += a*b
    yield abs(k) % 10

def stepMany(vals, n):
  for i in range(n):
    vals = list(step(vals))
  return vals

print stepMany(vals, 100)[:8]
