package main

import (
	"bufio"
	"fmt"
	"io"
	debug "log"
	"os"
)

// for temporary use only
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	answer, err := solve(os.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func solve(reader io.Reader) (answer string, err error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		v := scanner.Text()

		// modify stuff here

		debug.Println(v)

	}
	err = scanner.Err()
	if err != nil {
		return
	}

	// modify stuff here

	answer = "unimplemented"
	return
}
