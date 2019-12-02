package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func solve(lines []string) (interface{}, error) {
	for _, line := range lines {
		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}

		fmt.Printf("v=%#v\n", v)
	}

	answer := "unimplemented"
	return answer, nil
}

func main() {
	lines, err := getLines()
	if err != nil {
		panic(err)
	}

	answer, err := solve(lines)
	if err != nil {
		panic(err)
	}

	fmt.Println(answer)
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

var source = os.Stdin

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
