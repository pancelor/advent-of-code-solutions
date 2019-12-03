package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func solve(wires []Wire) (interface{}, error) {
	assert(len(wires) == 2, "wrong # wires")
	wireA := wires[0]
	wireB := wires[1]

	// fmt.Printf("wireA=%#v\n", wireA)
	// fmt.Printf("wireB=%#v\n", wireB)
	minDist := 1000000
	for _, lineA := range wireA.lines {
		for _, lineB := range wireB.lines {
			x, y := lineA.intersect(lineB)
			if x == 0 && y == 0 {
				continue
			}
			dist := abs(x) + abs(y)
			if dist < minDist {
				minDist = dist
			}
		}
	}

	return minDist, nil
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

// Wire .
type Wire struct {
	lines []Line
}

// Line .
type Line struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

func (l *Line) validate() error {
	if l.x1 != l.x2 && l.y1 != l.y2 {
		return errors.New("diagonal line")
	}
	return nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func bubble(a, b *int) {
	if *b < *a {
		temp := *b
		*b = *a
		*a = temp
	}
}

// OneDCollide does a 1-dimensional collision test
// returns whether [a1, a2) and [b1, b2) have any intersection
// auto-sorts if a2 < a1 etc
func OneDCollide(a1, a2, b1, b2 int) bool {
	bubble(&a1, &a2)
	bubble(&b1, &b2)
	if b2 <= a1 {
		return false
	}
	if a2 <= b1 {
		return false
	}
	return true
}

func between(x, a, b int) bool {
	// inclusive
	bubble(&a, &b)
	res := a <= x && x <= b
	// fmt.Println(a, x, b, res)
	return res
}

func (l *Line) intersect(o Line) (int, int) {
	s := 0
	if l.x1 == l.x2 {
		s += 0b10
	}
	if o.x1 == o.x2 {
		s += 0b01
	}
	switch s {
	case 0b11:
		// ||
		if l.x1 != o.x1 {
			return 0, 0
		}
		fmt.Println("||")
	case 0b10:
		// |-
		x := l.x1
		y := o.y1
		if !between(x, o.x1, o.x2) ||
			!between(y, l.y1, l.y2) {
			return 0, 0
		}
		return x, y
	case 0b01:
		// -|
		x := o.x1
		y := l.y1
		if !between(x, l.x1, l.x2) ||
			!between(y, o.y1, o.y2) {
			return 0, 0
		}
		return x, y
	case 0b00:
		// --
		if l.y1 != o.y1 {
			return 0, 0
		}
		fmt.Println("--")
	default:
		panic(fmt.Sprintf("impossible intersection %d", s))
	}
	return 0, 0
}

func test() {
	f := func(a, b Line, x, y int, msg string) {
		var xr, yr int
		// fmt.Printf("running %s (1)\n", msg)
		xr, yr = a.intersect(b)
		assert(x == xr && y == yr, fmt.Sprintf("expected (%d,%d) got (%d, %d)\n", x, y, xr, yr))
		// fmt.Printf("running %s (2)\n", msg)
		xr, yr = b.intersect(a)
		assert(x == xr && y == yr, fmt.Sprintf("expected (%d,%d) got (%d, %d)\n", x, y, xr, yr))
	}
	l1 := Line{x1: 0, y1: 0, x2: 0, y2: 10}
	l2 := Line{x1: 0, y1: 0, x2: 10, y2: 0}
	l3 := Line{x1: -4, y1: 6, x2: 4, y2: 6}
	l4 := Line{x1: 5, y1: 6, x2: 10, y2: 6}
	f(l1, l2, 0, 0, "t1")
	f(l1, l3, 0, 6, "t2")
	f(l2, l3, 0, 0, "t3")
	f(l1, l4, 0, 0, "t4")

	l5 := Line{x1: 3, x2: 3, y1: -5, y2: -2}
	l6 := Line{x1: 6, x2: 2, y1: -3, y2: -3}
	f(l5, l6, 3, -3, "t5")
}

func getInput() ([]Wire, error) {
	lines, err := getLines()
	if err != nil {
		return nil, err
	}

	var res []Wire
	for _, l := range lines {
		if l == "" {
			continue
		}
		var segs []Line
		var x, y int
		for _, token := range strings.Split(l, ",") {
			var dirCode string
			var dist int
			_, err := fmt.Sscanf(token, "%1s%d", &dirCode, &dist)
			if err != nil {
				return res, err
			}

			x1 := x
			y1 := y
			switch dirCode {
			case "R":
				x += dist
			case "U":
				y -= dist
			case "L":
				x -= dist
			case "D":
				y += dist
			}
			s := Line{x1: x1, y1: y1, x2: x, y2: y}
			if err := s.validate(); err != nil {
				return res, err
			}
			segs = append(segs, s)
		}
		res = append(res, Wire{lines: segs})
	}

	return res, nil
}

var source = os.Stdin

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
