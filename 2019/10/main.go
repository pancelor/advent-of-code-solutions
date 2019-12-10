package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
)

func inbounds(x, a, b int) bool {
	return a <= x && x < b
}

func bubble(x, y *int) {
	if *y < *x {
		temp := *y
		*y = *x
		*x = temp
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(x, y int) int {
	bubble(&x, &y)
	x = abs(x)
	y = abs(y)

	// fmt.Println(x, y)
	if x == 0 {
		// fmt.Println(y)
		return y
	}
	return gcd(x, y%x)
}

func solve(input *Input) (interface{}, error) {
	minBlocked := 1000
	aBest := Asteroid{x: 0, y: 0}
	for _, a1 := range input.list {
		// fmt.Printf("\nchecking (%d, %d)\n", a1.x, a1.y)

		blocked := make(map[Asteroid]bool)
		for _, a2 := range input.list {
			if a1 == a2 {
				continue
			}
			if blocked[a2] {
				// fmt.Printf("skipping (%d, %d)\n", a2.x, a2.y)
				continue
			}
			dx := a2.x - a1.x
			dy := a2.y - a1.y
			d := gcd(dx, dy)
			// fmt.Printf("gcd(%d,%d)=%d\n", dx, dy, d)
			dx /= d
			dy /= d

			x := a2.x
			y := a2.y
			for {
				x += dx
				y += dy
				if !(inbounds(x, 0, input.w) && inbounds(y, 0, input.h)) {
					break
				}
				a3 := input.grid[y][x]
				if a3 {
					// fmt.Printf("found (%d, %d)\n", x, y)
					blocked[Asteroid{x: x, y: y}] = true
				}
			}
		}
		nBlocked := len(blocked)
		if nBlocked < minBlocked {
			minBlocked = nBlocked
			aBest = a1
		}
		// fmt.Printf("A(%d,%d)=%d\n", a1.x, a1.y, nBlocked)
	}

	fmt.Printf("aBest=%#v\n", aBest)

	a1 := aBest
	depth := make(map[Asteroid]int)
	for _, a2 := range input.list {
		if a1 == a2 {
			continue
		}
		if depth[a2] != 0 {
			// fmt.Printf("skipping (%d, %d)\n", a2.x, a2.y)
			continue
		}
		depth[a2] = 0
		dx := a2.x - a1.x
		dy := a2.y - a1.y
		d := gcd(dx, dy)
		// fmt.Printf("gcd(%d,%d)=%d\n", dx, dy, d)
		dx /= d
		dy /= d

		x := a2.x
		y := a2.y
		n := 1
		for {
			x += dx
			y += dy
			if !(inbounds(x, 0, input.w) && inbounds(y, 0, input.h)) {
				break
			}
			a3 := input.grid[y][x]
			if a3 {
				// fmt.Printf("found (%d, %d)\n", x, y)
				depth[Asteroid{x: x, y: y}] = n
				n++
			}
		}
	}

	maxDepth := 0
	invertedDepths := invertMap(depth)
	for k := range invertedDepths {
		// fmt.Printf("|%#v: %#v\n", k, v)
		if maxDepth < k {
			maxDepth = k
		}
	}

	// fmt.Printf("depth=%#v\n", depth)
	// fmt.Printf("invertedDepths=%#v\n", invertedDepths)
	fmt.Printf("%d: %v\n", maxDepth, invertedDepths[maxDepth])

	vaporized := 0
	for i := 0; i < 10000; i++ {
		layer := invertedDepths[i]
		dV := len(layer)
		println("dV:", dV)
		if vaporized+dV > 200 {
			sort.Sort(ByAngle(layer))
			fmt.Printf("find #%d of %v\n", 200-vaporized, layer)
			return layer[200-vaporized], nil
		}
		vaporized += dV
	}
	return nil, fmt.Errorf("no answer")
}

// ByAngle .
type ByAngle []Asteroid

func (a ByAngle) Len() int      { return len(a) }
func (a ByAngle) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByAngle) Less(i, j int) bool {
	return clockAngle(a[i].x-11, -(a[i].y-19)) < clockAngle(a[j].x-11, -(a[j].y-19))
}

func clockAngle(x, y int) float64 {
	t := math.Atan2(float64(y), float64(x))
	if math.Pi/2 < t && t <= math.Pi {
		t -= 2 * math.Pi
	}
	return 90 - t*180/math.Pi
}

func invertMap(m map[Asteroid]int) map[int][]Asteroid {
	res := make(map[int][]Asteroid)
	for k, v := range m {
		arr, prs := res[v]
		if !prs {
			// unecessary?
			arr = make([]Asteroid, 0)
		}
		arr = append(arr, k)
		res[v] = arr
	}
	return res
}

func test() {
	assert(gcd(6, 10) == 2, "t1")
	assert(gcd(2, 10) == 2, "t2")
	assert(gcd(3, 5) == 1, "t3")

	noon := clockAngle(0, 1)
	q1 := clockAngle(1, 1)
	q4 := clockAngle(1, -2)
	q3 := clockAngle(-5, -2)
	q2 := clockAngle(-5, 8)
	fmt.Printf("%f | %f %f %f %f\n", noon, q1, q4, q3, q2)
	assert(noon < q1, "t4")
	assert(q1 < q4, "t5")
	assert(q4 < q3, "t6")
	assert(q3 < q2, "t7")

	// assert(false, "exit after tests")
}

func main() {
	test()

	input, err := getInput()
	if err != nil {
		panic(err)
	}

	answer, err := solve(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("answer:\n%v\n", answer)
}

//
// helpers
//

// for temporary use only
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// for temporary use only
func assert(b bool, msg string) {
	if !b {
		panic(errors.New(msg))
	}
}

var source = os.Stdin

// Asteroid .
type Asteroid struct {
	x int
	y int
}

// Input .
type Input struct {
	list []Asteroid
	grid [][]bool
	w    int
	h    int
}

func getInput() (*Input, error) {
	lines, err := getLines()
	if err != nil {
		return nil, err
	}

	var grid [][]bool
	var list []Asteroid
	for r, l := range lines {
		if l == "" {
			continue
		}
		var gridLine []bool
		for c, ch := range []byte(l) {
			var v bool
			switch ch {
			case '.':
				v = false
			case '#':
				v = true
			default:
				panic(fmt.Sprintf("bad input char '%s'", ch))
			}
			gridLine = append(gridLine, v)
			if v {
				list = append(list, Asteroid{x: c, y: r})
			}
		}
		grid = append(grid, gridLine)
	}

	res := Input{
		list: list,
		grid: grid,
		w:    len(grid[0]),
		h:    len(grid),
	}
	return &res, nil
}

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
