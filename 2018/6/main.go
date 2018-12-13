package main

import (
	"bytes"
	"errors"
	debug "log"
	"io"
	"fmt"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	data, err := getInput()
	if err != nil {
		panic(err)
	}
	answer, err := solve(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func getInput() ([]point, error) {
	out := make([]point, 0)
	for i := 0; true; i++{
		var x, y int
		_, err := fmt.Scanf("%d, %d", &x, &y)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		out = append(out, newPoint(x, y, byte('A'+i)))
	}
	return out, nil
}

func solve(data []point) (answer int, err error) {
	if len(data) == 0 {
		err = errors.New("no data")
		return
	}
	// debug.Println(data)
	maxX := max(fmapPoints(func(p point) int {return p.x}, data))
	maxY := max(fmapPoints(func(p point) int {return p.y}, data))
	minX := min(fmapPoints(func(p point) int {return p.x}, data))
	minY := min(fmapPoints(func(p point) int {return p.y}, data))
	// debug.Println(minX, minY, maxX, maxY)

	minX -= 2
	minY -= 2
	maxX += 2
	maxY += 2
	g := newGrid(minX, minY, maxX, maxY)
	for _, p := range data {
		err = g.setP(p)
		if err != nil {
			return
		}
	}
	// fmt.Println(g)

	distTo := func(r, c int, p point) int {
		return abs(r-p.y) + abs(c-p.x)
	}
	area := make([]int, len(data))
	infinite := make([]bool, len(data))
	for r := minY; r < maxY; r++ {
		for c := minX; c < maxX; c++ {
			closestPointId := 0
			closestDist := distTo(r, c, data[0])
			isContested := false
			for i, p := range data {
				if i == 0 {
					continue
				}
				dist := distTo(r, c, p)
				if dist < closestDist {
					closestPointId = i
					closestDist = dist
					isContested = false
				} else if dist == closestDist {
					isContested = true
				}
			}
			if closestDist != 0 {
				var marker byte
				if isContested {
					marker = '.'
				} else {
					area[closestPointId] += 1
					if r == minY || r == maxY - 1 || c == minX || c == maxX - 1 {
						infinite[closestPointId] = true
					}
					marker = byteToLower(data[closestPointId].id)
				}
				g.setP(newPoint(c, r, marker))
			}
		}
	}
	fmt.Println(g)
	maxArea := 0
	initd := false
	for i, p := range(data) {
		if !infinite[i] && (!initd || area[i] > maxArea) {
			initd = true
			maxArea = area[i]
		}
		fmt.Printf("%c: %d (inf: %v). max is %d\n", p.id, area[i], infinite[i], maxArea)
	}
	fmt.Println()

	answer = maxArea
	return
}

func byteToLower(b byte) byte {
	return bytes.ToLower([]byte{b})[0]
}

func pointsAtDistanceFrom(d int, center point) []point {
	res := make([]point, 0)
	id := byteToLower(center.id)
	for dx := 0; dx < d; dx++ {
		dy := d - dx
		res = append(res,
			newPoint(center.x + dx, center.y + dy, id),
			newPoint(center.x - dy, center.y + dx, id),
			newPoint(center.x - dx, center.y - dy, id),
			newPoint(center.x + dy, center.y - dx, id),
		)
	}
	return res
}



type point struct {
	x, y int
	id byte
}

func newPoint(x, y int, id byte) point {
	return point {
		x: x,
		y: y,
		id: id,
	}
}

func (p point) String() string {
	return fmt.Sprintf("Point(%c@%d,%d)", p.id, p.y, p.x) //row-col
}

func fmapPoints(f func(point) int, arr []point) []int {
	res := make([]int, len(arr))
	for i, x := range arr {
		res[i] = f(x)
	}
	return res
}



type grid struct {
	g [][]string //gType
	xOffset, yOffset int
	width, height int
}

func newGrid(minX, minY, maxX, maxY int) grid {
	if maxX < minX || maxY < minY{
		panic("bad grid")
	}
	width := maxX-minX
	height := maxY-minY
	g := make([][]string, height) //gType
	for r := 0; r < height; r++ {
		g[r] = make([]string, width) //gType
		for c := 0; c < width; c++ {
			g[r][c] = " " // gType
		}
	}
	return grid{
		g: g,
		xOffset: -minX,
		yOffset: -minY,
		width: width,
		height: height,
	}
}

func (g grid) get(r, c int) string { //gType
	return g.g[r+g.yOffset][c+g.xOffset]
}

func (g grid) set(r, c int, val string) (err error) { //gType
	rr := r+g.yOffset
	cc := c+g.xOffset
	if !(0 <= rr && rr < g.height && 0 <= cc && cc < g.width) {
		err = errors.New("out of bounds")
		return
	}
	g.g[rr][cc] = val
	return
}

func (g grid) setP(p point) error {
	return g.set(p.y, p.x, fmt.Sprintf("%c", p.id)) //gType
}

func (g grid) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "[")
	h := g.height
	if h > 0 {
		w := g.width
		for r := 0; r < h; r++ {
			if r != 0 {
				fmt.Fprintf(&b, "\n ")
			}
			fmt.Fprintf(&b, "[")
			for c := 0; c < w; c++ {
				if c != 0 {
					fmt.Fprintf(&b, " ")
				}
				fmt.Fprintf(&b, "%s", g.g[r][c])
			}
			fmt.Fprintf(&b, "]")
		}
	}
	fmt.Fprintf(&b, "]\n")
	return b.String()
}


func max(arr []int) int {
	if len(arr) == 0 {
		panic("max([]int{}) is undefined")
	}
	res := arr[0]
	for _, x := range arr {
		if x > res {
			res = x
		}
	}
	return res
}

func min(arr []int) int {
	if len(arr) == 0 {
		panic("min([]int{}) is undefined")
	}
	res := arr[0]
	for _, x := range arr {
		if x < res {
			res = x
		}
	}
	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

func argmin(arr []int) int {
	if len(arr) == 0 {
		panic("argmin([]int{}) is undefined")
	}
	imin := 0
	min := arr[0]
	for i := 1; i < len(arr); i++ {
		elem := arr[i]
		if elem < min {
			imin = i
			min = elem
		}
	}
	return imin
}
