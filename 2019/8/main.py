#!/usr/bin/env python

from collections import Counter

full = raw_input()

def getLayers():
  i = 0
  while i*25*6 < len(full):
    layer = full[i*25*6:(i+1)*25*6]
    yield layer
    i += 1

def count(layer, target):
  return Counter(layer)[target]

l = min(getLayers(), key=lambda l: count(l, '0'))
print count(l, '1')*count(l, '2')
