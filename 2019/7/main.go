package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/computer"
	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

// Input is the type of input we pass to solve()
type Input []int

func solve(memTemplate Input) interface{} {
	best := 0
	for _, settings := range allPerms {
		assert(len(settings) == 5, "bad settings len")

		var cpus []computer.CPU
		var lastOut chan int
		for i, name := range []string{"A", "B", "C", "D", "E"} {
			cpu := computer.MakeCPU(name)
			cpus = append(cpus, cpu)
			cpu.SetMemory(memTemplate)
			if i != 0 {
				cpu.InChan = lastOut
			}
			lastOut = cpu.OutChan
			cpu.Run()
			cpu.InChan <- settings[i]
		}

		// go
		cpus[0].InChan <- 0

		go func() {
			for x := range cpus[4].OutChan {
				// fmt.Printf("---got %d from cpus[4].OutChan---\n", x)
				cpus[0].InChan <- x
			}
		}()

		<-cpus[4].DoneChan
		res := <-cpus[0].InChan
		fmt.Printf("(settings=%v) res=%v\n", settings, res)
		if res > best {
			best = res
		}
	}

	return best
}

func init() {
	// tests go here
	assert(true, "t1")

	// assert(false, "exit after tests")
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

var allPerms = [][5]int{
	[5]int{0, 1, 2, 3, 4},
	[5]int{0, 1, 2, 4, 3},
	[5]int{0, 1, 3, 2, 4},
	[5]int{0, 1, 3, 4, 2},
	[5]int{0, 1, 4, 2, 3},
	[5]int{0, 1, 4, 3, 2},
	[5]int{0, 2, 1, 3, 4},
	[5]int{0, 2, 1, 4, 3},
	[5]int{0, 2, 3, 1, 4},
	[5]int{0, 2, 3, 4, 1},
	[5]int{0, 2, 4, 1, 3},
	[5]int{0, 2, 4, 3, 1},
	[5]int{0, 3, 1, 2, 4},
	[5]int{0, 3, 1, 4, 2},
	[5]int{0, 3, 2, 1, 4},
	[5]int{0, 3, 2, 4, 1},
	[5]int{0, 3, 4, 1, 2},
	[5]int{0, 3, 4, 2, 1},
	[5]int{0, 4, 1, 2, 3},
	[5]int{0, 4, 1, 3, 2},
	[5]int{0, 4, 2, 1, 3},
	[5]int{0, 4, 2, 3, 1},
	[5]int{0, 4, 3, 1, 2},
	[5]int{0, 4, 3, 2, 1},
	[5]int{1, 0, 2, 3, 4},
	[5]int{1, 0, 2, 4, 3},
	[5]int{1, 0, 3, 2, 4},
	[5]int{1, 0, 3, 4, 2},
	[5]int{1, 0, 4, 2, 3},
	[5]int{1, 0, 4, 3, 2},
	[5]int{1, 2, 0, 3, 4},
	[5]int{1, 2, 0, 4, 3},
	[5]int{1, 2, 3, 0, 4},
	[5]int{1, 2, 3, 4, 0},
	[5]int{1, 2, 4, 0, 3},
	[5]int{1, 2, 4, 3, 0},
	[5]int{1, 3, 0, 2, 4},
	[5]int{1, 3, 0, 4, 2},
	[5]int{1, 3, 2, 0, 4},
	[5]int{1, 3, 2, 4, 0},
	[5]int{1, 3, 4, 0, 2},
	[5]int{1, 3, 4, 2, 0},
	[5]int{1, 4, 0, 2, 3},
	[5]int{1, 4, 0, 3, 2},
	[5]int{1, 4, 2, 0, 3},
	[5]int{1, 4, 2, 3, 0},
	[5]int{1, 4, 3, 0, 2},
	[5]int{1, 4, 3, 2, 0},
	[5]int{2, 0, 1, 3, 4},
	[5]int{2, 0, 1, 4, 3},
	[5]int{2, 0, 3, 1, 4},
	[5]int{2, 0, 3, 4, 1},
	[5]int{2, 0, 4, 1, 3},
	[5]int{2, 0, 4, 3, 1},
	[5]int{2, 1, 0, 3, 4},
	[5]int{2, 1, 0, 4, 3},
	[5]int{2, 1, 3, 0, 4},
	[5]int{2, 1, 3, 4, 0},
	[5]int{2, 1, 4, 0, 3},
	[5]int{2, 1, 4, 3, 0},
	[5]int{2, 3, 0, 1, 4},
	[5]int{2, 3, 0, 4, 1},
	[5]int{2, 3, 1, 0, 4},
	[5]int{2, 3, 1, 4, 0},
	[5]int{2, 3, 4, 0, 1},
	[5]int{2, 3, 4, 1, 0},
	[5]int{2, 4, 0, 1, 3},
	[5]int{2, 4, 0, 3, 1},
	[5]int{2, 4, 1, 0, 3},
	[5]int{2, 4, 1, 3, 0},
	[5]int{2, 4, 3, 0, 1},
	[5]int{2, 4, 3, 1, 0},
	[5]int{3, 0, 1, 2, 4},
	[5]int{3, 0, 1, 4, 2},
	[5]int{3, 0, 2, 1, 4},
	[5]int{3, 0, 2, 4, 1},
	[5]int{3, 0, 4, 1, 2},
	[5]int{3, 0, 4, 2, 1},
	[5]int{3, 1, 0, 2, 4},
	[5]int{3, 1, 0, 4, 2},
	[5]int{3, 1, 2, 0, 4},
	[5]int{3, 1, 2, 4, 0},
	[5]int{3, 1, 4, 0, 2},
	[5]int{3, 1, 4, 2, 0},
	[5]int{3, 2, 0, 1, 4},
	[5]int{3, 2, 0, 4, 1},
	[5]int{3, 2, 1, 0, 4},
	[5]int{3, 2, 1, 4, 0},
	[5]int{3, 2, 4, 0, 1},
	[5]int{3, 2, 4, 1, 0},
	[5]int{3, 4, 0, 1, 2},
	[5]int{3, 4, 0, 2, 1},
	[5]int{3, 4, 1, 0, 2},
	[5]int{3, 4, 1, 2, 0},
	[5]int{3, 4, 2, 0, 1},
	[5]int{3, 4, 2, 1, 0},
	[5]int{4, 0, 1, 2, 3},
	[5]int{4, 0, 1, 3, 2},
	[5]int{4, 0, 2, 1, 3},
	[5]int{4, 0, 2, 3, 1},
	[5]int{4, 0, 3, 1, 2},
	[5]int{4, 0, 3, 2, 1},
	[5]int{4, 1, 0, 2, 3},
	[5]int{4, 1, 0, 3, 2},
	[5]int{4, 1, 2, 0, 3},
	[5]int{4, 1, 2, 3, 0},
	[5]int{4, 1, 3, 0, 2},
	[5]int{4, 1, 3, 2, 0},
	[5]int{4, 2, 0, 1, 3},
	[5]int{4, 2, 0, 3, 1},
	[5]int{4, 2, 1, 0, 3},
	[5]int{4, 2, 1, 3, 0},
	[5]int{4, 2, 3, 0, 1},
	[5]int{4, 2, 3, 1, 0},
	[5]int{4, 3, 0, 1, 2},
	[5]int{4, 3, 0, 2, 1},
	[5]int{4, 3, 1, 0, 2},
	[5]int{4, 3, 1, 2, 0},
	[5]int{4, 3, 2, 0, 1},
	[5]int{4, 3, 2, 1, 0},
}
