package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func solve(input []int) (interface{}, error) {
	for _, val := range input {
		fmt.Printf("val=%#v\n", val)
	}

	answer := "unimplemented"
	return answer, nil
}

func test() {
	assert(true, "t1")
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

func getInput() ([]int, error) {
	lines, err := getLines()
	if err != nil {
		return nil, err
	}

	var res []int
	for _, l := range lines {
		if l == "" {
			continue
		}

		v, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}

		res = append(res, v)
	}

	return res, nil
}

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
