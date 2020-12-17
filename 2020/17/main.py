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
from copy import deepcopy

class Sim:
	def __init__(self,grid):
		self.chunks=[]
		self.n=0
		self.w=len(grid[0])
		self.h=len(grid)
		for y,row in enumerate(grid):
			for x,v in enumerate(row):
				self.set(x,y,0,0,v)
	def chunk_for(self,x,y,z,w):
		chunk=None
		for c in self.chunks:
			if c["w"]==w and c["z"]==z and c["x"]==x/8 and c["y"]==y/8:
				chunk=c
		if chunk==None:
			grid=[]
			for _ in range(8):
				grid.append([False]*8)
			chunk={"z":z,"w":w,"y":y/8,"x":x/8,"grid":grid}
			self.chunks.append(chunk)
			# print "new chunk!",len(self.chunks)
		return chunk
	def get(self,x,y,z,w):
		chunk=self.chunk_for(x,y,z,w)
		return chunk["grid"][y%8][x%8]
	def set(self,x,y,z,w,val):
		chunk=self.chunk_for(x,y,z,w)
		chunk["grid"][y%8][x%8]=val
	def step(self):
		old=deepcopy(self)
		self.n+=1
		for w in range(-self.n,self.n+1):
			print "w",w
			for z in range(-self.n,self.n+1):
				print " z",z
				for y in range(-self.n,self.h+self.n+1):
					for x in range(-self.n,self.w+self.n+1):
						count=0
						for nx,ny,nz,nw in neighbors(x,y,z,w):
							if old.get(nx,ny,nz,nw):
								count+=1
						if old.get(x,y,z,w):
							if not (count==2 or count==3):
								self.set(x,y,z,w,False)
						else:
							if count==3:
								self.set(x,y,z,w,True)
	def count(self):
		c=0
		for w in range(-self.n,self.n+1):
			for z in range(-self.n,self.n+1):
				for y in range(-self.n,self.h+self.n+1):
					for x in range(-self.n,self.w+self.n+1):
						if self.get(x,y,z,w):
							c+=1
		return c
	def __str__(self):
		s=""
		for c in self.chunks:
			s+="x,y,z,w=%d,%d,%d,%d\n"%(c["x"]*8,c["y"]*8,c["z"],c["w"])
			s+=self.chunkstr(c)
		return s
	def chunkstr(self,c):
		s=""
		for row in c["grid"]:
			s+=''.join("#" if x else "." for x in row)+"\n"
		return s

def neighbors(x,y,z,w):
	for dw in range(-1,2):
		for dz in range(-1,2):
			for dy in range(-1,2):
				for dx in range(-1,2):
					if dx==0 and dy==0 and dz==0 and dw==0:
						continue
					yield x+dx,y+dy,z+dz,w+dw

grid=[]
for line in sys.stdin:
	row=[c=="#" for c in line.strip()]
	grid.append(row)

sim=Sim(grid)
# print sim

for i in range(6):
	print "STEP"
	sim.step()
	# print sim
print sim.count()
