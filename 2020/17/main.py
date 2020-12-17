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
	def __init__(self,):
		self.chunks=[]
	def chunk_for(self,x,y,z):
		chunk=None
		for c in self.chunks:
			if c.z==z and c.x==x/8 and c.y==y/8:
				chunk=c
		if chunk==None:
			grid=[]
			for _ in range(8):
				grid.append([False]*8)
			chunk={"z":z,"y":y/8,"x":x/8,"grid":grid}
			self.chunks.append(chunk)
			print "new chunk!",chunk
		return chunk
	def get(self,x,y,z):
		chunk=self.chunk_for(x,y,z)
		return chunk[y%8][x%8]
	def set(self,x,y,z,val):
		chunk=self.chunk_for(x,y,z)
		chunk[y%8][x%8]=val
	def __str__(self):
		for c in self.chunks:
			print "x,y,z",c["x"]*8,c["y"]*8,c["z"]
			for row in c.grid:
				print ''.join("#" if x else "." for x in row)

def neighbors(x,y,z):
	for dz in range(-1,2):
		for dy in range(-1,2):
			for dx in range(-1,2):
				if dx==0 and dy==0 and dz==0:
					continue
				yield x+dx,y+dy,z+dz

def make_empty(slices,n):
	s=[]
	for y,row in enumerate(slices[0]):
		row2=[]
		for x,val in enumerate(row):
			row2.append(False)
		s.append(row2)
	slices[n]=s

def get(slices,x,y,z):
	if z in slices:
		if 0<=x and x<len(slices[0][0]):
			if 0<=y and y<len(slices[0]):
				return slices[z][y][x]
	return False

def step(slices):
	newslices=deepcopy(slices)
	n=slices["n"]
	n+=1
	newslices["n"]=n
	make_empty(newslices,-n)
	make_empty(newslices,n)
	for z in range(-n,n+1):
		for y in range(len(slices[0])):
			for x in range(len(slices[0][0])):
				count=0
				# print "x,y,z",x,y,z
				for nx,ny,nz in neighbors(x,y,z):
					if get(slices,nx,ny,nz):
						# print "  nx,ny,nz",nx,ny,nz
						count+=1
				if get(slices,x,y,z):
					if not (count==2 or count==3):
						newslices[z][y][x]=False
				else:
					if count==3:
						newslices[z][y][x]=True
	return newslices

slices={"n":0}
slices[0]=[]
for line in sys.stdin:
	row=[c=="#" for c in line.strip()]
	slices[0].append(row)

pslice(slices[0])
# for y in range(3):
# 	for x in range(3):
# 		print get(slices,x,y,0)

for i in range(6):
	slices=step(slices)

count=0
n=slices["n"]
for z in range(-n,n+1):
	print z
	pslice(slices[z])
	for y,row in enumerate(slices[z]):
		for x,val in enumerate(row):
			if val:
				count+=1
print count
