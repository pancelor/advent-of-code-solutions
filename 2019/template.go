package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

// Input is the type of input we pass to solve()
type Input []int

func solve(in Input) interface{} {
	for _, val := range in {
		fmt.Printf("val=%#v\n", val)
	}

	answer := "unimplemented"
	return answer
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
