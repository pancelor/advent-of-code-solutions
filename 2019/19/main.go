package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/computer"
	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

type TileType int

const (
	TT_UNKNOWN TileType = iota
	TT_UNEXPLORED
	TT_EXPLORED
	TT_WALL
)

func (t *TileType) String() string {
	switch *t {
	case TT_UNKNOWN:
		return "-"
	case TT_UNEXPLORED:
		return "*"
	case TT_EXPLORED:
		return "."
	case TT_WALL:
		return "#"
	}
	return "?"
}

type point struct {
	x int
	y int
}

func oppDir(dir int) int {
	return (dir + 2) % 4
}

func (p point) addDir(dir int, n int) point {
	newP := p
	switch dir {
	case 0:
		newP.x += n
	case 1:
		newP.y -= n
	case 2:
		newP.x -= n
	case 3:
		newP.y += n
	default:
		assert(false, fmt.Sprintf("bad dir %d", dir))
	}
	return newP
}

func dirString(dir int) string {
	switch dir {
	case 0:
		return ">"
	case 1:
		return "^"
	case 2:
		return "<"
	case 3:
		return "v"
	}
	return "?"
}

type Grid map[point]TileType

func makeGrid() Grid {
	return make(map[point]TileType)
}

func (g *Grid) bounds() (int, int, int, int) {
	var minX, maxX, minY, maxY int
	for p := range *g {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	assert(minX <= maxX, "bad bounds")
	assert(minY <= maxY, "bad bounds")
	return minX, maxX, minY, maxY
}

func (g Grid) String() string {
	var b strings.Builder
	x1, x2, y1, y2 := g.bounds()
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			t := g[point{x, y}]
			fmt.Fprintf(&b, t.String())
		}
		fmt.Fprintf(&b, "\n")
	}
	return b.String()
}

func (p *point) nz() bool {
	return *p != (point{})
}

func translate(dir int) int {
	switch dir {
	case 0:
		return 4
	case 1:
		return 1
	case 2:
		return 3
	case 3:
		return 2
	}
	return 0
}

type robot struct {
	cpu     *computer.CPU
	pos     point
	dir     int
	grid    Grid
	goalPos point
}

// record tile type seen at current pos
func (r *robot) record(t TileType) {
	// fmt.Printf("Recording %s at %v\n", t.String(), r.pos)
	existing := r.grid[r.pos]
	if existing == TT_UNKNOWN || existing == TT_UNEXPLORED {
		r.grid[r.pos] = t
	}
}

func (r *robot) Draw() string {
	var b strings.Builder
	x1, x2, y1, y2 := r.grid.bounds()
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			p := point{x, y}
			if p == r.goalPos {
				fmt.Fprint(&b, "$")
			} else if p == r.pos {
				fmt.Fprint(&b, dirString(r.dir))
			} else {
				t := r.grid[p]
				fmt.Fprintf(&b, t.String())
			}
		}
		fmt.Fprintf(&b, "\n")
	}
	return b.String()
}

func (r *robot) step(dist int) {
	r.pos = r.pos.addDir(r.dir, dist)
}

// takes a normal dir (0123) and returns if was able to move in that direction
func (r *robot) move(dir int, real bool) bool {
	{ // input
		r.dir = dir
		r.cpu.SendInput(translate(dir))
		r.step(1)
	}
	{ // output
		code := r.cpu.RecvOutput()

		switch code {
		case 0:
			// fmt.Println("Hit a wall")
			r.record(TT_WALL)
			r.step(-1)
			return false
		case 1:
			// fmt.Println("Success")
			if real {
				r.record(TT_EXPLORED)
			} else {
				r.record(TT_UNEXPLORED)
			}
			return true
		case 2:
			// fmt.Println("Success; found goal!")
			if real {
				r.record(TT_EXPLORED)
			} else {
				r.record(TT_UNEXPLORED)
			}
			r.goalPos = r.pos

			assert(r.goalPos.nz(), "goal is at 0,0")
			// fmt.Printf("found goal! %v\n", r.goalPos)

			return true
		default:
			assert(false, "unknown return code %d", code)
		}
	}
	return false
}

func (r *robot) unexploredNeighbors() []int {
	var res []int
	for d := 0; d < 4; d++ {
		t := r.grid[r.pos.addDir(d, 1)]
		if t == TT_UNEXPLORED {
			res = append(res, d)
		}
	}
	return res
}

func makeRobot(cpu *computer.CPU) *robot {
	r := robot{
		cpu:  cpu,
		grid: makeGrid(),
	}
	r.record(TT_EXPLORED)
	return &r
}

func promptDir() int {
	fmt.Printf("dir=? ")
	code := helpers.ReadLine()
	switch code {
	case "d":
		return 0
	case "w":
		return 1
	case "a":
		return 2
	case "s":
		return 3
	default:
		fmt.Printf("(wasd only)\n")
		return promptDir()
	}
	return 0
}

func (r *robot) look() {
	temp := r.dir
	for d := 0; d < 4; d++ {
		if r.move(d, false) {
			assert(r.move(oppDir(d), false), "couldn't unmove")
		}
	}
	r.dir = temp
}

// uses backtracking to explore the entire area
// assumes the maze is a DAG; does not check to make sure this is true
// fullExplore: if true, explore EVERYTHING. if false, move until on top of end
// returns the number of moves
func (r *robot) explore(fullExplore bool) int {
	start := r.pos
	var hist []int // history of directions moved
	maxDepth := 0
	r.look()
	for {
		unexplored := r.unexploredNeighbors()
		var dir int
		amRewinding := false
		if len(unexplored) == 0 {
			if r.pos == start {
				// we've backtracked all the way back to the start
				return maxDepth
			}
			var last int
			last, hist = hist[len(hist)-1], hist[:len(hist)-1] // pop
			dir = oppDir(last)
			amRewinding = true
		} else {
			amRewinding = false
			dir = unexplored[0]
		}

		success := r.move(dir, true)
		if success && !amRewinding {
			hist = append(hist, dir)
		}
		if !success && amRewinding {
			assert(false, "couldnt rewind")
		}
		r.look()

		depth := len(hist)
		if depth > maxDepth {
			maxDepth = depth
		}

		if fullExplore {
			// fmt.Printf("depth=%v, maxDepth=%v\n%s\n", depth, maxDepth, r.Draw())
			// time.Sleep(3 * time.Millisecond)
		}
		if r.pos == r.goalPos && !fullExplore {
			return depth
		}
	}
}

var reNum = regexp.MustCompile("^\\d+$")

func stringSeq(s string) (res []int) {
	for _, word := range strings.Split(s, ",") {
		var toSend int
		var err error
		if reNum.MatchString(word) {
			toSend, err = strconv.Atoi(word)
			helpers.Check(err)
		} else {
			toSend = int([]byte(word)[0])
		}
		res = append(res, toSend)
		res = append(res, int(','))
	}
	res = append(res, int('\n'))
	return
}

func test(in []int, p point) {
	for y := p.y - 1; y < p.y+105; y++ {
		for x := p.x - 1; x < p.x+105; x++ {
			if getPos(in, x, y) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func solve(in []int) interface{} {
	p := solveFRD(in)
	test(in, p)
	return p.x*10000 + p.y
}

func solveFRD(in []int) point {
	Width := 100
	Start := 2017
	for d := Start; true; d++ {
		fmt.Println(d)
		streak := 0
	INNER:
		for i := d; i >= 0; i-- {
			y := i
			x := d - i
			if getPos(in, x, y) {
				streak++
				if streak == Width {
					return point{x: x - (Width - 1), y: y}
				}
			} else {
				if streak > 0 {
					fmt.Println("nope")
					break INNER
				}
			}
		}
	}
	return point{}
}

func getPos(in []int, x int, y int) bool {
	cpu := computer.MakeCPU("sal")
	cpu.SetMemory(in)
	cpu.Run()

	cpu.SendInput(x)
	cpu.SendInput(y)
	return cpu.RecvOutput() == 1
}

func main() {
	answer := solve(input)
	fmt.Printf("answer:\n%v\n", answer)
}

var input []int = []int{
	109, 424, 203, 1, 21101, 11, 0, 0, 1106, 0, 282, 21102, 18, 1, 0, 1105, 1, 259, 1201, 1, 0, 221, 203, 1, 21101, 0, 31, 0, 1105, 1, 282, 21102, 1, 38, 0, 1106, 0, 259, 21001, 23, 0, 2, 22101, 0, 1, 3, 21102, 1, 1, 1, 21101, 57, 0, 0, 1106, 0, 303, 1202, 1, 1, 222, 21002, 221, 1, 3, 20102, 1, 221, 2, 21101, 259, 0, 1, 21102, 1, 80, 0, 1105, 1, 225, 21101, 83, 0, 2, 21101, 91, 0, 0, 1105, 1, 303, 2102, 1, 1, 223, 20102, 1, 222, 4, 21102, 259, 1, 3, 21101, 225, 0, 2, 21101, 225, 0, 1, 21102, 1, 118, 0, 1105, 1, 225, 20102, 1, 222, 3, 21101, 0, 51, 2, 21102, 1, 133, 0, 1105, 1, 303, 21202, 1, -1, 1, 22001, 223, 1, 1, 21102, 1, 148, 0, 1106, 0, 259, 1201, 1, 0, 223, 21002, 221, 1, 4, 21002, 222, 1, 3, 21101, 13, 0, 2, 1001, 132, -2, 224, 1002, 224, 2, 224, 1001, 224, 3, 224, 1002, 132, -1, 132, 1, 224, 132, 224, 21001, 224, 1, 1, 21102, 195, 1, 0, 106, 0, 108, 20207, 1, 223, 2, 21002, 23, 1, 1, 21102, -1, 1, 3, 21101, 0, 214, 0, 1106, 0, 303, 22101, 1, 1, 1, 204, 1, 99, 0, 0, 0, 0, 109, 5, 2102, 1, -4, 249, 21202, -3, 1, 1, 21202, -2, 1, 2, 22102, 1, -1, 3, 21101, 0, 250, 0, 1105, 1, 225, 22102, 1, 1, -4, 109, -5, 2106, 0, 0, 109, 3, 22107, 0, -2, -1, 21202, -1, 2, -1, 21201, -1, -1, -1, 22202, -1, -2, -2, 109, -3, 2106, 0, 0, 109, 3, 21207, -2, 0, -1, 1206, -1, 294, 104, 0, 99, 22101, 0, -2, -2, 109, -3, 2105, 1, 0, 109, 5, 22207, -3, -4, -1, 1206, -1, 346, 22201, -4, -3, -4, 21202, -3, -1, -1, 22201, -4, -1, 2, 21202, 2, -1, -1, 22201, -4, -1, 1, 22101, 0, -2, 3, 21101, 0, 343, 0, 1106, 0, 303, 1105, 1, 415, 22207, -2, -3, -1, 1206, -1, 387, 22201, -3, -2, -3, 21202, -2, -1, -1, 22201, -3, -1, 3, 21202, 3, -1, -1, 22201, -3, -1, 2, 21202, -4, 1, 1, 21102, 384, 1, 0, 1106, 0, 303, 1105, 1, 415, 21202, -4, -1, -4, 22201, -4, -3, -4, 22202, -3, -2, -2, 22202, -2, -4, -4, 22202, -3, -2, -3, 21202, -4, -1, -2, 22201, -3, -2, 1, 22101, 0, 1, -4, 109, -5, 2105, 1, 0,
}
