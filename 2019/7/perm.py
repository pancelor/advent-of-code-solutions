#!/usr/bin/env python

import itertools as itt

print 'var allPerms = [][5]int{'
for a,b,c,d,e in itt.permutations([0,1,2,3,4]):
  print("  [5]int{}{}, {}, {}, {}, {}{},".format("{",a,b,c,d,e,"}"))
print '}'
