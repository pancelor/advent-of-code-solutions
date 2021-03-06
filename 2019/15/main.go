package main

import (
	"fmt"
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

func solve(in []int) interface{} {
	cpu := computer.MakeCPU("robbie")
	cpu.SetMemory(in)
	// fmt.Println(cpu.PrintProgram())
	cpu.Run()

	r := makeRobot(&cpu)
	// return r.explore(false)
	r.explore(false)
	assert(r.pos == r.goalPos, "robot not on goal")
	fmt.Printf("%s\n", r.Draw())

	r.grid = makeGrid()
	return r.explore(true)
}

func main() {
	answer := solve(prog)
	fmt.Printf("answer:\n%v\n", answer)
}

var prog = []int{
	3, 1033, 1008, 1033, 1, 1032, 1005, 1032, 31, 1008, 1033, 2, 1032, 1005, 1032, 58, 1008, 1033, 3, 1032, 1005, 1032, 81, 1008, 1033, 4, 1032, 1005, 1032, 104, 99, 101, 0, 1034, 1039, 1001, 1036, 0, 1041, 1001, 1035, -1, 1040, 1008, 1038, 0, 1043, 102, -1, 1043, 1032, 1, 1037, 1032, 1042, 1105, 1, 124, 101, 0, 1034, 1039, 101, 0, 1036, 1041, 1001, 1035, 1, 1040, 1008, 1038, 0, 1043, 1, 1037, 1038, 1042, 1106, 0, 124, 1001, 1034, -1, 1039, 1008, 1036, 0, 1041, 1001, 1035, 0, 1040, 1001, 1038, 0, 1043, 1002, 1037, 1, 1042, 1105, 1, 124, 1001, 1034, 1, 1039, 1008, 1036, 0, 1041, 102, 1, 1035, 1040, 101, 0, 1038, 1043, 1002, 1037, 1, 1042, 1006, 1039, 217, 1006, 1040, 217, 1008, 1039, 40, 1032, 1005, 1032, 217, 1008, 1040, 40, 1032, 1005, 1032, 217, 1008, 1039, 1, 1032, 1006, 1032, 165, 1008, 1040, 33, 1032, 1006, 1032, 165, 1101, 0, 2, 1044, 1106, 0, 224, 2, 1041, 1043, 1032, 1006, 1032, 179, 1101, 1, 0, 1044, 1106, 0, 224, 1, 1041, 1043, 1032, 1006, 1032, 217, 1, 1042, 1043, 1032, 1001, 1032, -1, 1032, 1002, 1032, 39, 1032, 1, 1032, 1039, 1032, 101, -1, 1032, 1032, 101, 252, 1032, 211, 1007, 0, 43, 1044, 1105, 1, 224, 1101, 0, 0, 1044, 1106, 0, 224, 1006, 1044, 247, 1002, 1039, 1, 1034, 1002, 1040, 1, 1035, 102, 1, 1041, 1036, 1001, 1043, 0, 1038, 101, 0, 1042, 1037, 4, 1044, 1105, 1, 0, 13, 30, 60, 64, 5, 28, 36, 24, 67, 12, 1, 67, 32, 39, 14, 78, 29, 17, 38, 88, 79, 9, 62, 25, 15, 18, 88, 25, 7, 81, 38, 41, 10, 69, 86, 32, 11, 33, 1, 10, 22, 84, 14, 92, 48, 79, 10, 3, 62, 33, 61, 13, 93, 78, 20, 63, 68, 17, 80, 34, 12, 8, 23, 61, 90, 51, 17, 84, 37, 46, 64, 25, 3, 73, 19, 45, 99, 41, 62, 21, 77, 8, 17, 89, 9, 13, 84, 75, 85, 14, 53, 60, 6, 29, 76, 63, 14, 23, 63, 61, 93, 72, 17, 41, 28, 94, 5, 3, 19, 47, 57, 55, 14, 34, 38, 79, 85, 40, 13, 22, 99, 67, 72, 15, 62, 15, 6, 63, 3, 90, 2, 87, 20, 84, 15, 50, 70, 27, 18, 78, 21, 70, 48, 52, 2, 99, 92, 55, 3, 46, 41, 93, 99, 88, 13, 39, 4, 45, 71, 3, 96, 1, 91, 59, 31, 53, 23, 25, 82, 32, 50, 16, 60, 38, 78, 34, 59, 30, 15, 51, 92, 3, 22, 26, 62, 60, 37, 42, 74, 28, 21, 76, 7, 24, 70, 18, 40, 11, 81, 41, 9, 73, 62, 12, 66, 81, 9, 3, 74, 62, 11, 6, 56, 16, 34, 20, 78, 79, 1, 97, 17, 39, 87, 15, 12, 77, 94, 28, 22, 66, 45, 59, 39, 2, 6, 52, 6, 72, 49, 17, 92, 15, 86, 18, 92, 79, 67, 20, 22, 72, 10, 72, 3, 52, 26, 77, 78, 41, 97, 36, 59, 88, 24, 57, 12, 38, 90, 53, 14, 38, 67, 2, 36, 44, 93, 99, 10, 41, 49, 3, 16, 7, 63, 32, 11, 15, 81, 12, 91, 39, 62, 19, 83, 6, 91, 28, 19, 80, 38, 23, 63, 31, 71, 14, 58, 8, 21, 71, 21, 21, 81, 38, 26, 32, 29, 82, 52, 28, 72, 54, 97, 41, 65, 96, 75, 1, 48, 28, 80, 66, 25, 47, 49, 29, 87, 51, 12, 50, 70, 36, 60, 81, 29, 77, 76, 55, 25, 40, 45, 83, 91, 26, 72, 99, 12, 47, 11, 20, 27, 52, 9, 98, 17, 99, 27, 37, 62, 25, 3, 15, 73, 66, 22, 5, 85, 5, 20, 98, 20, 38, 62, 78, 21, 16, 59, 28, 98, 38, 31, 2, 40, 46, 87, 14, 48, 33, 80, 48, 36, 27, 56, 21, 1, 50, 83, 3, 61, 92, 20, 52, 16, 50, 10, 86, 9, 98, 39, 56, 25, 50, 42, 39, 91, 81, 56, 25, 70, 44, 24, 15, 99, 4, 20, 55, 12, 98, 27, 65, 20, 77, 97, 76, 36, 42, 87, 6, 11, 79, 65, 16, 65, 44, 13, 90, 13, 48, 79, 13, 95, 60, 19, 55, 24, 66, 4, 53, 11, 23, 68, 14, 97, 53, 45, 14, 16, 93, 18, 29, 83, 5, 6, 77, 19, 70, 97, 34, 20, 70, 52, 11, 74, 14, 72, 10, 36, 44, 33, 45, 19, 38, 36, 77, 5, 37, 51, 1, 55, 17, 2, 48, 23, 18, 2, 34, 90, 97, 24, 30, 51, 66, 33, 70, 51, 37, 31, 51, 37, 65, 55, 18, 8, 66, 4, 65, 62, 26, 93, 29, 88, 3, 75, 73, 24, 23, 67, 1, 13, 68, 7, 36, 87, 62, 48, 1, 31, 45, 28, 62, 86, 24, 98, 1, 59, 49, 37, 26, 62, 36, 44, 66, 18, 17, 97, 92, 40, 36, 65, 80, 84, 5, 84, 6, 79, 87, 36, 31, 96, 15, 71, 96, 2, 72, 11, 81, 95, 94, 41, 54, 31, 58, 25, 74, 24, 51, 81, 38, 32, 73, 22, 96, 40, 62, 22, 59, 74, 39, 25, 86, 2, 55, 20, 61, 40, 37, 88, 69, 1, 60, 42, 18, 31, 54, 13, 27, 19, 93, 34, 41, 99, 33, 89, 20, 16, 52, 84, 32, 94, 31, 6, 61, 25, 1, 61, 1, 38, 78, 87, 39, 31, 39, 26, 68, 42, 36, 2, 94, 66, 2, 67, 30, 80, 2, 95, 65, 40, 54, 50, 33, 11, 23, 97, 89, 1, 31, 56, 9, 35, 49, 92, 55, 23, 84, 48, 91, 20, 7, 72, 25, 55, 3, 85, 3, 16, 40, 90, 22, 99, 44, 38, 86, 98, 11, 76, 26, 76, 13, 82, 80, 24, 93, 4, 15, 64, 95, 58, 15, 85, 25, 57, 29, 66, 3, 66, 19, 98, 57, 24, 44, 59, 35, 76, 48, 31, 92, 33, 94, 68, 56, 41, 45, 15, 46, 5, 68, 15, 65, 34, 73, 49, 68, 17, 78, 28, 80, 24, 59, 26, 74, 21, 52, 1, 94, 5, 61, 41, 88, 37, 56, 1, 49, 0, 0, 21, 21, 1, 10, 1, 0, 0, 0, 0, 0, 0,
}
