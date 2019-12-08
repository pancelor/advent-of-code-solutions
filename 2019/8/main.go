package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func count(l Layer, i byte) int {
	num := 0
	for _, row := range l.arr {
		for _, x := range row {
			if x == i {
				num++
			}
		}
	}
	return num
}

func solve(layers []Layer) (interface{}, error) {
	argMin := -1
	min := 0
	for i, l := range layers {
		numZeros := count(l, '0')
		if argMin == -1 || numZeros < min {
			argMin = i
			min = numZeros
		}
		fmt.Printf("i, numZeros, argmMin, min=%d,%d,%d,%d\n", i, numZeros, argMin, min)
	}

	return count(layers[argMin], '2') * count(layers[argMin], '1'), nil
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

// Layer .
type Layer struct {
	arr [H][W]byte
}

// W width
const W = 25

// H height
const H = 6

func getInput() ([]Layer, error) {
	lines, err := getLines()
	if err != nil {
		return nil, err
	}

	var layers []Layer
	full := []byte(strings.Join(lines, ""))
	// fmt.Printf("full=%#v\n", full)
	for i := 0; i*W*H < len(full); i++ {
		l := [H][W]byte{}
		for r := 0; r < H; r++ {
			row := [W]byte{}
			for c := 0; c < W; c++ {
				row[c] = full[i*W*H+r*W+c]
			}
			l[r] = row
		}
		layers = append(layers, Layer{l})
	}
	return layers, nil
}

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
