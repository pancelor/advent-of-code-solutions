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
	debug.Println(data)
	maxX := max(fmapPoints(func(p point) int {return p.x}, data))
	maxY := max(fmapPoints(func(p point) int {return p.y}, data))
	minX := -max(fmapPoints(func(p point) int {return -p.x}, data))
	minY := -max(fmapPoints(func(p point) int {return -p.y}, data))
	// debug.Println(minX, minY, maxX, maxY)

	g := newGrid(minX-2, minY-2, maxX+2, maxY+2)
	for _, p := range data {
		err = g.setP(p)
		if err != nil {
			return
		}
	}
	fmt.Println(g)


// CAPITAL(point)
//   actually same as SOLID?
// NEUTRAL()
//   equally close to two points
// SOLID(point)
//   solidified distance from a point
// TENATIVE(point)
//   tenative distance from a point (to know whether it needs to be changed to NEUTRAL)
// UNEXPLORED()

// for i = 0; true; i++
//   for each CAPITAL p
//     ns = neighbors at distance i from p
//     for each ns n
//       switch n.type
//       CAPITAL: pass // e.g if two capitals are touching
//       NEUTRAL: pass
//       SOLID: pass
//       TENATIVE: n = NEUTRAL
//       UNEXPLORED: n = TENATIVE(p)

	for i := 0; i < 2; i++ {
		for _, p := range(data) {
			neighbors := pointsAtDistanceFrom(i+1, p)
			debug.Println(i, neighbors)
			for _, n := range(neighbors) {
				g.setP(n)
			}
		}
		fmt.Println(g)
	}

	answer = 0
	return
}

func pointsAtDistanceFrom(d int, center point) []point {
	res := make([]point, 0)
	id := bytes.ToLower([]byte{center.id})[0]
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
	return fmt.Sprintf("Point(%c: %d, %d)", p.id, p.x, p.y)
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
			g[r][c] = "   " // gType
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
	return g.set(p.y, p.x, fmt.Sprintf(" %c ", p.id)) //gType
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