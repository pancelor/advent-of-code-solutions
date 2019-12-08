package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func count(l Layer, i int) int {
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
	// 0 transparent, 1 not, 2 not
	res := Layer{}
	for _, l := range layers {
		for r, row := range l.arr {
			for c, x := range row {
				if x != 0 {
					res.arr[r][c] = x
				}
			}
		}
	}

	for _, row := range res.arr {
		for _, x := range row {
			switch x {
			case 2:
				fmt.Printf("%d", x)
			default:
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}

	return nil, nil
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
	arr [H][W]int
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
		l := [H][W]int{}
		for r := 0; r < H; r++ {
			row := [W]int{}
			for c := 0; c < W; c++ {
				ch := full[i*W*H+r*W+c]
				switch ch {
				case '0':
					row[c] = 2
				case '1':
					row[c] = 1
				}
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
