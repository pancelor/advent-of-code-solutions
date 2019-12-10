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

func solve(memTemplate Input) interface{} {
	cpu := computer.MakeCPU("A")
	cpu.SetMemory(memTemplate)
	cpu.Run()
	cpu.InChan <- 2

	var res []int
	for {
		select {
		case <-cpu.DoneChan:
			return res
		case x := <-cpu.OutChan:
			res = append(res, x)
		}
	}

	return res
}

func main() {
	input, err := getInput()
	check(err)
	answer := solve(input)
	fmt.Printf("answer:\n%v\n", answer)
}

// Input is the type of input we pass to solve()
type Input []int

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
