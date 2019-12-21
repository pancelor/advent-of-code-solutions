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

var inQueue []byte

func promptInput() {
	fmt.Printf("> ")
	str := helpers.ReadLine() + "\n"
	// fmt.Println("READ:", str)
	inQueue = []byte(str)
}

func solve(in []int) interface{} {
	cpu := computer.MakeCPU("sal")
	cpu.SetMemory(in)
	// fmt.Println(cpu.PrintProgram())
	cpu.Run()

	var lastOut int
	for {
		state := <-cpu.StateChan
		switch state {
		case computer.CS_WAITING_INPUT:
			if len(inQueue) == 0 {
				promptInput()
			}
			var b byte
			b, inQueue = inQueue[0], inQueue[1:]
			// fmt.Printf("sending %q\n", b)
			cpu.InChan <- int(b)
		case computer.CS_WAITING_OUTPUT:
			out := <-cpu.OutChan
			lastOut = out
			fmt.Print(string(out))
		case computer.CS_DONE:
			return lastOut
		}
	}
}

func main() {
	answer := solve(input)
	fmt.Printf("answer:\n%v\n", answer)
}

var input = []int{
	109, 2050, 21101, 966, 0, 1, 21102, 1, 13, 0, 1105, 1, 1378, 21102, 20, 1, 0, 1106, 0, 1337, 21101, 0, 27, 0, 1106, 0, 1279, 1208, 1, 65, 748, 1005, 748, 73, 1208, 1, 79, 748, 1005, 748, 110, 1208, 1, 78, 748, 1005, 748, 132, 1208, 1, 87, 748, 1005, 748, 169, 1208, 1, 82, 748, 1005, 748, 239, 21102, 1041, 1, 1, 21102, 73, 1, 0, 1105, 1, 1421, 21101, 78, 0, 1, 21101, 1041, 0, 2, 21102, 88, 1, 0, 1106, 0, 1301, 21102, 1, 68, 1, 21101, 1041, 0, 2, 21102, 1, 103, 0, 1105, 1, 1301, 1101, 0, 1, 750, 1106, 0, 298, 21102, 82, 1, 1, 21101, 1041, 0, 2, 21101, 125, 0, 0, 1106, 0, 1301, 1101, 2, 0, 750, 1106, 0, 298, 21102, 79, 1, 1, 21101, 0, 1041, 2, 21101, 147, 0, 0, 1105, 1, 1301, 21101, 84, 0, 1, 21102, 1041, 1, 2, 21102, 162, 1, 0, 1106, 0, 1301, 1101, 3, 0, 750, 1106, 0, 298, 21102, 65, 1, 1, 21102, 1, 1041, 2, 21102, 184, 1, 0, 1105, 1, 1301, 21102, 1, 76, 1, 21102, 1, 1041, 2, 21101, 0, 199, 0, 1106, 0, 1301, 21102, 75, 1, 1, 21101, 1041, 0, 2, 21102, 1, 214, 0, 1106, 0, 1301, 21101, 0, 221, 0, 1106, 0, 1337, 21101, 0, 10, 1, 21102, 1, 1041, 2, 21102, 236, 1, 0, 1105, 1, 1301, 1106, 0, 553, 21102, 1, 85, 1, 21101, 1041, 0, 2, 21101, 0, 254, 0, 1105, 1, 1301, 21102, 78, 1, 1, 21101, 0, 1041, 2, 21102, 1, 269, 0, 1105, 1, 1301, 21101, 276, 0, 0, 1106, 0, 1337, 21101, 0, 10, 1, 21101, 0, 1041, 2, 21102, 1, 291, 0, 1105, 1, 1301, 1101, 0, 1, 755, 1105, 1, 553, 21101, 0, 32, 1, 21101, 1041, 0, 2, 21101, 0, 313, 0, 1105, 1, 1301, 21102, 1, 320, 0, 1106, 0, 1337, 21102, 1, 327, 0, 1106, 0, 1279, 2102, 1, 1, 749, 21101, 0, 65, 2, 21102, 1, 73, 3, 21101, 0, 346, 0, 1105, 1, 1889, 1206, 1, 367, 1007, 749, 69, 748, 1005, 748, 360, 1101, 0, 1, 756, 1001, 749, -64, 751, 1106, 0, 406, 1008, 749, 74, 748, 1006, 748, 381, 1102, -1, 1, 751, 1106, 0, 406, 1008, 749, 84, 748, 1006, 748, 395, 1101, -2, 0, 751, 1105, 1, 406, 21101, 1100, 0, 1, 21102, 406, 1, 0, 1105, 1, 1421, 21101, 0, 32, 1, 21102, 1100, 1, 2, 21102, 1, 421, 0, 1105, 1, 1301, 21102, 1, 428, 0, 1105, 1, 1337, 21102, 1, 435, 0, 1105, 1, 1279, 2101, 0, 1, 749, 1008, 749, 74, 748, 1006, 748, 453, 1101, 0, -1, 752, 1106, 0, 478, 1008, 749, 84, 748, 1006, 748, 467, 1101, -2, 0, 752, 1106, 0, 478, 21102, 1168, 1, 1, 21101, 0, 478, 0, 1105, 1, 1421, 21101, 485, 0, 0, 1105, 1, 1337, 21102, 10, 1, 1, 21101, 1168, 0, 2, 21101, 0, 500, 0, 1105, 1, 1301, 1007, 920, 15, 748, 1005, 748, 518, 21101, 0, 1209, 1, 21102, 518, 1, 0, 1106, 0, 1421, 1002, 920, 3, 529, 1001, 529, 921, 529, 1002, 750, 1, 0, 1001, 529, 1, 537, 1001, 751, 0, 0, 1001, 537, 1, 545, 101, 0, 752, 0, 1001, 920, 1, 920, 1106, 0, 13, 1005, 755, 577, 1006, 756, 570, 21101, 1100, 0, 1, 21102, 1, 570, 0, 1105, 1, 1421, 21102, 1, 987, 1, 1106, 0, 581, 21101, 1001, 0, 1, 21101, 588, 0, 0, 1105, 1, 1378, 1102, 1, 758, 594, 101, 0, 0, 753, 1006, 753, 654, 21001, 753, 0, 1, 21102, 610, 1, 0, 1106, 0, 667, 21101, 0, 0, 1, 21101, 621, 0, 0, 1105, 1, 1463, 1205, 1, 647, 21102, 1015, 1, 1, 21101, 635, 0, 0, 1105, 1, 1378, 21102, 1, 1, 1, 21102, 1, 646, 0, 1105, 1, 1463, 99, 1001, 594, 1, 594, 1106, 0, 592, 1006, 755, 664, 1101, 0, 0, 755, 1106, 0, 647, 4, 754, 99, 109, 2, 1101, 726, 0, 757, 21201, -1, 0, 1, 21101, 9, 0, 2, 21101, 697, 0, 3, 21102, 1, 692, 0, 1106, 0, 1913, 109, -2, 2106, 0, 0, 109, 2, 101, 0, 757, 706, 1201, -1, 0, 0, 1001, 757, 1, 757, 109, -2, 2106, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 63, 223, 127, 159, 191, 95, 0, 183, 170, 218, 117, 217, 86, 155, 244, 187, 178, 56, 109, 253, 154, 38, 184, 213, 110, 189, 233, 212, 87, 50, 141, 169, 125, 236, 207, 171, 175, 230, 114, 137, 197, 99, 49, 249, 136, 239, 93, 242, 57, 231, 100, 228, 172, 219, 102, 140, 69, 84, 206, 92, 174, 77, 179, 166, 76, 163, 139, 203, 120, 173, 43, 62, 202, 124, 216, 103, 138, 71, 39, 54, 143, 227, 79, 157, 158, 167, 250, 126, 198, 235, 58, 121, 251, 156, 199, 35, 68, 246, 215, 53, 47, 185, 221, 182, 46, 168, 115, 51, 123, 142, 229, 85, 111, 201, 248, 107, 204, 70, 98, 34, 42, 188, 222, 60, 254, 196, 162, 220, 153, 78, 61, 55, 119, 252, 186, 181, 243, 238, 101, 118, 106, 214, 234, 226, 113, 108, 94, 152, 116, 200, 232, 59, 190, 247, 122, 245, 205, 237, 241, 177, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 73, 110, 112, 117, 116, 32, 105, 110, 115, 116, 114, 117, 99, 116, 105, 111, 110, 115, 58, 10, 13, 10, 87, 97, 108, 107, 105, 110, 103, 46, 46, 46, 10, 10, 13, 10, 82, 117, 110, 110, 105, 110, 103, 46, 46, 46, 10, 10, 25, 10, 68, 105, 100, 110, 39, 116, 32, 109, 97, 107, 101, 32, 105, 116, 32, 97, 99, 114, 111, 115, 115, 58, 10, 10, 58, 73, 110, 118, 97, 108, 105, 100, 32, 111, 112, 101, 114, 97, 116, 105, 111, 110, 59, 32, 101, 120, 112, 101, 99, 116, 101, 100, 32, 115, 111, 109, 101, 116, 104, 105, 110, 103, 32, 108, 105, 107, 101, 32, 65, 78, 68, 44, 32, 79, 82, 44, 32, 111, 114, 32, 78, 79, 84, 67, 73, 110, 118, 97, 108, 105, 100, 32, 102, 105, 114, 115, 116, 32, 97, 114, 103, 117, 109, 101, 110, 116, 59, 32, 101, 120, 112, 101, 99, 116, 101, 100, 32, 115, 111, 109, 101, 116, 104, 105, 110, 103, 32, 108, 105, 107, 101, 32, 65, 44, 32, 66, 44, 32, 67, 44, 32, 68, 44, 32, 74, 44, 32, 111, 114, 32, 84, 40, 73, 110, 118, 97, 108, 105, 100, 32, 115, 101, 99, 111, 110, 100, 32, 97, 114, 103, 117, 109, 101, 110, 116, 59, 32, 101, 120, 112, 101, 99, 116, 101, 100, 32, 74, 32, 111, 114, 32, 84, 52, 79, 117, 116, 32, 111, 102, 32, 109, 101, 109, 111, 114, 121, 59, 32, 97, 116, 32, 109, 111, 115, 116, 32, 49, 53, 32, 105, 110, 115, 116, 114, 117, 99, 116, 105, 111, 110, 115, 32, 99, 97, 110, 32, 98, 101, 32, 115, 116, 111, 114, 101, 100, 0, 109, 1, 1005, 1262, 1270, 3, 1262, 21002, 1262, 1, 0, 109, -1, 2105, 1, 0, 109, 1, 21102, 1, 1288, 0, 1105, 1, 1263, 21001, 1262, 0, 0, 1102, 1, 0, 1262, 109, -1, 2105, 1, 0, 109, 5, 21101, 0, 1310, 0, 1106, 0, 1279, 22101, 0, 1, -2, 22208, -2, -4, -1, 1205, -1, 1332, 22102, 1, -3, 1, 21102, 1332, 1, 0, 1105, 1, 1421, 109, -5, 2105, 1, 0, 109, 2, 21101, 1346, 0, 0, 1105, 1, 1263, 21208, 1, 32, -1, 1205, -1, 1363, 21208, 1, 9, -1, 1205, -1, 1363, 1106, 0, 1373, 21101, 0, 1370, 0, 1105, 1, 1279, 1106, 0, 1339, 109, -2, 2105, 1, 0, 109, 5, 2102, 1, -4, 1385, 21001, 0, 0, -2, 22101, 1, -4, -4, 21101, 0, 0, -3, 22208, -3, -2, -1, 1205, -1, 1416, 2201, -4, -3, 1408, 4, 0, 21201, -3, 1, -3, 1106, 0, 1396, 109, -5, 2105, 1, 0, 109, 2, 104, 10, 22102, 1, -1, 1, 21102, 1436, 1, 0, 1106, 0, 1378, 104, 10, 99, 109, -2, 2105, 1, 0, 109, 3, 20002, 594, 753, -1, 22202, -1, -2, -1, 201, -1, 754, 754, 109, -3, 2106, 0, 0, 109, 10, 21101, 0, 5, -5, 21102, 1, 1, -4, 21102, 0, 1, -3, 1206, -9, 1555, 21101, 3, 0, -6, 21102, 1, 5, -7, 22208, -7, -5, -8, 1206, -8, 1507, 22208, -6, -4, -8, 1206, -8, 1507, 104, 64, 1106, 0, 1529, 1205, -6, 1527, 1201, -7, 716, 1515, 21002, 0, -11, -8, 21201, -8, 46, -8, 204, -8, 1106, 0, 1529, 104, 46, 21201, -7, 1, -7, 21207, -7, 22, -8, 1205, -8, 1488, 104, 10, 21201, -6, -1, -6, 21207, -6, 0, -8, 1206, -8, 1484, 104, 10, 21207, -4, 1, -8, 1206, -8, 1569, 21101, 0, 0, -9, 1105, 1, 1689, 21208, -5, 21, -8, 1206, -8, 1583, 21102, 1, 1, -9, 1105, 1, 1689, 1201, -5, 716, 1589, 20102, 1, 0, -2, 21208, -4, 1, -1, 22202, -2, -1, -1, 1205, -2, 1613, 21201, -5, 0, 1, 21101, 0, 1613, 0, 1106, 0, 1444, 1206, -1, 1634, 22101, 0, -5, 1, 21101, 1627, 0, 0, 1105, 1, 1694, 1206, 1, 1634, 21101, 0, 2, -3, 22107, 1, -4, -8, 22201, -1, -8, -8, 1206, -8, 1649, 21201, -5, 1, -5, 1206, -3, 1663, 21201, -3, -1, -3, 21201, -4, 1, -4, 1106, 0, 1667, 21201, -4, -1, -4, 21208, -4, 0, -1, 1201, -5, 716, 1676, 22002, 0, -1, -1, 1206, -1, 1686, 21101, 1, 0, -4, 1106, 0, 1477, 109, -10, 2106, 0, 0, 109, 11, 21101, 0, 0, -6, 21102, 0, 1, -8, 21101, 0, 0, -7, 20208, -6, 920, -9, 1205, -9, 1880, 21202, -6, 3, -9, 1201, -9, 921, 1725, 20101, 0, 0, -5, 1001, 1725, 1, 1732, 21002, 0, 1, -4, 21202, -4, 1, 1, 21102, 1, 1, 2, 21101, 9, 0, 3, 21102, 1, 1754, 0, 1105, 1, 1889, 1206, 1, 1772, 2201, -10, -4, 1766, 1001, 1766, 716, 1766, 21001, 0, 0, -3, 1105, 1, 1790, 21208, -4, -1, -9, 1206, -9, 1786, 22102, 1, -8, -3, 1105, 1, 1790, 21202, -7, 1, -3, 1001, 1732, 1, 1795, 21002, 0, 1, -2, 21208, -2, -1, -9, 1206, -9, 1812, 21201, -8, 0, -1, 1106, 0, 1816, 21202, -7, 1, -1, 21208, -5, 1, -9, 1205, -9, 1837, 21208, -5, 2, -9, 1205, -9, 1844, 21208, -3, 0, -1, 1106, 0, 1855, 22202, -3, -1, -1, 1106, 0, 1855, 22201, -3, -1, -1, 22107, 0, -1, -1, 1106, 0, 1855, 21208, -2, -1, -9, 1206, -9, 1869, 22101, 0, -1, -8, 1106, 0, 1873, 21201, -1, 0, -7, 21201, -6, 1, -6, 1106, 0, 1708, 22101, 0, -8, -10, 109, -11, 2106, 0, 0, 109, 7, 22207, -6, -5, -3, 22207, -4, -6, -2, 22201, -3, -2, -1, 21208, -1, 0, -6, 109, -7, 2106, 0, 0, 0, 109, 5, 2101, 0, -2, 1912, 21207, -4, 0, -1, 1206, -1, 1930, 21101, 0, 0, -4, 21201, -4, 0, 1, 21202, -3, 1, 2, 21101, 0, 1, 3, 21102, 1, 1949, 0, 1106, 0, 1954, 109, -5, 2105, 1, 0, 109, 6, 21207, -4, 1, -1, 1206, -1, 1977, 22207, -5, -3, -1, 1206, -1, 1977, 22101, 0, -5, -5, 1105, 1, 2045, 21201, -5, 0, 1, 21201, -4, -1, 2, 21202, -3, 2, 3, 21101, 1996, 0, 0, 1106, 0, 1954, 22101, 0, 1, -5, 21102, 1, 1, -2, 22207, -5, -3, -1, 1206, -1, 2015, 21102, 1, 0, -2, 22202, -3, -2, -3, 22107, 0, -4, -1, 1206, -1, 2037, 21201, -2, 0, 1, 21102, 2037, 1, 0, 106, 0, 1912, 21202, -3, -1, -3, 22201, -5, -3, -5, 109, -6, 2105, 1, 0,
}
