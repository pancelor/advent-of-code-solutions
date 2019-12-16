package main

import (
	"fmt"
	"strconv"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

func solve(in []int) interface{} {
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

func getInput() ([]int, error) {
	lines, err := helpers.GetLines()
	if err != nil {
		return nil, err
	}

	var res []int
	for _, l := range []byte(lines[0]) {
		d := string(l)
		if d == "" {
			continue
		}

		v, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}

		res = append(res, v)
	}

	return res, nil
}
