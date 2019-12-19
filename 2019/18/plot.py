#!/usr/bin/env python

import matplotlib.pyplot as plt

vals = []
with open("num-keys.txt", 'r') as f:
  for line in f:
    vals.append(int(line))

plt.plot(vals)
plt.xlabel('iteration #')
plt.ylabel('# keys')
plt.show()
