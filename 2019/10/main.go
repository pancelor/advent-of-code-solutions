package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func inbounds(x, a, b int) bool {
	return a <= x && x < b
}

func solve(input *Input) (interface{}, error) {
	minBlocked := 1000
	for _, a1 := range input.list {
		blocked := make(map[Asteroid]bool)
		for _, a2 := range input.list {
			if a1 == a2 {
				continue
			}
			if blocked[a2] {
				// fmt.Printf("skipping blocked asteroid at %d %d\n", a2.x, a2.y)
				continue
			}
			dx := a2.x - a1.x
			dy := a2.y - a1.y
			x := a1.x
			y := a1.y
			for {
				x += dx
				y += dy
				if !(inbounds(x, 0, input.w) && inbounds(y, 0, input.h)) {
					break
				}
				a3 := input.grid[y][x]
				if a3 {
					// fmt.Printf("found blocked asteroid at %d %d\n", x, y)
					blocked[Asteroid{x: x, y: y}] = true
				}
			}
		}
		nBlocked := len(blocked)
		if nBlocked < minBlocked {
			minBlocked = nBlocked
		}
	}

	return len(input.list) - minBlocked, nil
}

func test() {
	assert(true, "t1")

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
			list = append(list, Asteroid{x: c, y: r})
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
