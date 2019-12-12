package helpers

import (
	"bufio"
	"fmt"
	"os"
)

// Check panics on non-nil errors
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Assert panics with msg unless b is true
func Assert(b bool, msg string, errorfArgs ...interface{}) {
	if !b {
		panic(fmt.Errorf(msg, errorfArgs...))
	}
}

// GetLines collects all lines from stdin in an array
func GetLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// ReadLine returns one line of stdin
func ReadLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	Assert(scanner.Err() == nil, "reading standard input: %s", scanner.Err())
	return ""
}

// EnsureInbounds makes sure the pointers won't overflow the buffer
func EnsureInbounds(mem []int, ptr ...int) {
	for _, p := range ptr {
		Assert(0 <= p && p < len(mem), "oob: %d", p)
	}
}
