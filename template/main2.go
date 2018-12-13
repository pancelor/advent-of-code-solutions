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

func getInput(in io.Reader) ([]prereq, error) {
	data := make([]prereq, 0)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		req := prereq{}
		fmt.Sscanf(line, "Step %c must be finished before step %c can begin.", &req.a, &req.b)
		data = append(data, req)
	}
	return data, scanner.Err()
}

func solve(data []prereq) (answer int, err error) {
	for _, v := range data {
		debug.Println(v)
	}

	return
}

// "<a> must be finished before <b> can begin."
type prereq struct{
	a, b byte
}

func (p prereq) String() string {
	return fmt.Sprintf("%c->%c", p.a, p.b)
}
