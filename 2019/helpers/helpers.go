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

// ExpectEqual panics unless a == b
func ExpectEqual(a, b interface{}) {
	if a != b {
		panic(fmt.Errorf("Expected %#v, got %#v", b, a))
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
// note: this is sorta fucky with piped input; just use a scanner directly
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

// Abs returns the absolute value
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Inbounds .
func Inbounds(x, a, b int) bool {
	return a <= x && x < b
}

// Bubble swaps the inputs if they're out of order
func Bubble(x, y *int) {
	if *y < *x {
		temp := *y
		*y = *x
		*x = temp
	}
}

// Gcd returns the least common multiple
func Gcd(xs ...int) int {
	res := 1
	for _, x := range xs {
		res = Gcd2(res, x)
	}
	return res
}

// Gcd2 returns the greatest common denominator of x and y
func Gcd2(x, y int) int {
	Bubble(&x, &y)
	x = Abs(x)
	y = Abs(y)

	// fmt.Println(x, y)
	if x == 0 {
		// fmt.Println(y)
		return y
	}
	return Gcd2(x, y%x)
}

func init() {
	ExpectEqual(Gcd2(6, 10), 2)
	ExpectEqual(Gcd2(6, 6), 6)
	ExpectEqual(Gcd2(5, 6), 1)
	ExpectEqual(Gcd2(1, 1), 1)
	ExpectEqual(Gcd2(1, 4), 1)
	ExpectEqual(Gcd2(0, 7), 7)
}

// Lcm returns the least common multiple
func Lcm(xs ...int) int {
	res := 1
	for _, x := range xs {
		res = Lcm2(res, x)
	}
	return res
}

// Lcm2 returns the least common multiple of x and y
func Lcm2(x, y int) int {
	return x * y / Gcd2(x, y)
}

func init() {
	ExpectEqual(Lcm2(6, 10), 30)
	ExpectEqual(Lcm2(6, 6), 6)
	ExpectEqual(Lcm2(5, 6), 30)
	ExpectEqual(Lcm2(1, 1), 1)
	ExpectEqual(Lcm2(1, 4), 4)
	ExpectEqual(Lcm2(0, 7), 0)
}
