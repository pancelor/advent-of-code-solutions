package main

import (
	"bufio"
	"fmt"
	"io"
	debug "log"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	init, rules, err := getInput(os.Stdin)
	check(err)
	answer, err := solve(init, rules)
	check(err)
	fmt.Println(answer)
}

type todo int

type state todo

type rule [6]bool

func plantFromBool(b bool) byte {
	if b {
		return '#'
	} else {
		return '.'
	}
}

func (r rule) String() string {
	var s strings.Builder
	for i := 0; i < 5; i++ {
		fmt.Fprintf(&s, "%c", plantFromBool(r[i]))
	}

	fmt.Fprintf(&s, " => %c", plantFromBool(r[5]))

	return s.String()
}

func getInput(in io.Reader) (init state, rules []rule, err error) {
	rules = make([]rule, 0)
	scanner := bufio.NewScanner(in)
	started := false
	for scanner.Scan() {
		line := scanner.Text()
		if !started {
			started = true
			var rest string
			fmt.Sscanf(line, "initial state: %s", &rest)
			// TODO build init
			debug.Println("rest", rest)
		} else {
			r := rule{}
			// TODO build r
			rules = append(rules, r)
		}
	}
	err = scanner.Err()
	return
}

func solve(init state, rules []rule) (answer int, err error) {
	debug.Println(init)
	for _, v := range rules {
		debug.Println(v)
	}

	return
}
