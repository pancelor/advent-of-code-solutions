package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
	"github.com/pancelor/advent-of-code-solutions/2019/computer"
)

var assert = helpers.Assert
var check = helpers.Check

// Input is the type of input we pass to solve()
type Input []int

type point struct {
	x int
	y int
}

func (p point) addDir(dir int, n int) point {
	newP := p
	switch (dir) {
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

func changeDir(d int, dd int) int {
	var dd2 int
	switch (dd) {
	case 0:
		dd2 = 1
	case 1:
		dd2 = -1
	default:
		assert(false, fmt.Sprintf("bad ddir %d", dd))
	}
	return (4+d+dd2) % 4
}
func init() {
	assert(changeDir(0, 0) == 1, "t1")
	assert(changeDir(0, 1) == 3, "t2")
	assert(changeDir(2, 0) == 3, "t3")
	assert(changeDir(2, 1) == 1, "t4")
}

type color int

func (c color) String() string {
	switch c {
	case 0:
		return "."
	case 1:
		return "#"
	}
	return "?"
}

func dirString(dir int) string {
	switch (dir) {
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

func draw(colors map[point]color, pos point, dir int) {
	for y := -5; y < 10; y++ {
		for x := -5; x < 55; x++ {
			p := point{x:x, y:y}
			if p == pos {
				fmt.Printf("%s", dirString(dir))
			} else {
				c := colors[p]
				fmt.Printf("%s", c.String())
			}
		}
		fmt.Printf("\n")
	}
}

func solve(prog Input) interface{} {
	colors := make(map[point]color)

	cpu := computer.MakeCPU("daniel")
	cpu.SetMemory(prog)
	cpu.Run()
	// cpu := makeFakeCPU([]int{1,0,0,0,1,0,1,0})
	// cpu := makeFakeCPU([]int{1,0,0,0,1,0,1,0,0,1,1,0,1,0})
	pos := point{0,0}
	dir := 1
	colors[pos] = color(1)
	for !cpu.Halted {
		currentColor := colors[pos]
		cpu.InChan <- int(currentColor)

		select {
		case colorCode := <- cpu.OutChan:
			colors[pos] = color(colorCode)

			dDir := <- cpu.OutChan
			dir = changeDir(dir, dDir)
			pos = pos.addDir(dir, 1)
		case <-cpu.DoneChan:
			fmt.Printf("done\n")
		}
	}

	draw(colors, pos, dir)

	return len(colors)
}

type fakeCPU struct {
	InChan chan int
	OutChan chan int
	DoneChan chan struct{}
	Halted bool
}

func makeFakeCPU(outputs []int) *fakeCPU {
	fake := fakeCPU{
		InChan: make(chan int),
		OutChan: make(chan int),
		DoneChan: make(chan struct{}),
	}

	go func() {
		for _ = range fake.InChan {}
	}()

	go func() {
		for _, x := range outputs {
			fake.OutChan <- x
		}
		fake.Halted = true
		fake.DoneChan <- struct{}{}
	}()

	return &fake
}

func main() {
	input, err := getInput()
	check(err)
	answer := solve(input)
	fmt.Printf("answer:\n%v\n", answer)
}

func getInput() (Input, error) {
	lines, err := helpers.GetLines()
	if err != nil {
		return nil, err
	}

	var res []int
	for _, l := range strings.Split(lines[0], ",") {
		if l == "" {
			continue
		}

		v, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}

		// var dirCode string
		// var dist int
		// _, err := fmt.Sscanf(token, "%1s%d", &dirCode, &dist)
		// if err != nil {
		// 	return res, err
		// }

		res = append(res, v)
	}

	return res, nil
}
