package main

import (
	"bufio"
	debug "log"
	"io"
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	data, err := getInput(os.Stdin)
	if err != nil {
		panic(err)
	}
	answer, err := solve(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func getInput(in io.Reader) ([]point, error) {
	out := make([]point, 0)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		var x, y int
		fmt.Sscanf(line, "%d, %d", &x, &y)
		out = append(out, newPoint(x, y))
	}
	return out, scanner.Err()
}

func solve(data []point) (answer string, outErr error) {
	for v := range data {
		debug.Printf("%#v", v)
	}
	answer = "unimplemented"
	return
}

type point struct {
	x, y int
}

func newPoint(x, y int) point {
	return point {
		x: x,
		y: y,
	}
}
