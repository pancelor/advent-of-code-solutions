package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func solve(low, high int) (interface{}, error) {
	fmt.Printf("low=%#v\n", low)
	fmt.Printf("high=%#v\n", high)

	var count int
	for i := low; i <= high; i++ {
		if works(strconv.Itoa(i)) {
			count++
		}
	}
	return count, nil
}

func works(s string) bool {
	hasDouble := false
	currentLen := 1
	last := byte('0')
	for _, ch := range []byte(s) {
		if ch < last {
			return false
		}
		if last == ch {
			currentLen++
		} else {
			if currentLen == 2 {
				hasDouble = true
			}
			currentLen = 1
		}
		last = ch
	}
	return hasDouble
}

func main() {
	low, high, err := getInput()
	if err != nil {
		panic(err)
	}

	answer, err := solve(low, high)
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

func getInput() (int, int, error) {
	lines, err := getLines()
	if err != nil {
		return 0, 0, err
	}
	assert(len(lines) > 0, "bad len lines")
	arr := strings.Split(lines[0], "-")
	assert(len(arr) == 2, "bad len arr")
	a, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.Atoi(arr[1])
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
