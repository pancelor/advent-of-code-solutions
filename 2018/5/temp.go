package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(killPair('a', "foobar"))
	fmt.Println(killPair('a', "foobAr"))
	fmt.Println(killPair('a', "foobAAaAr"))
	fmt.Println(killPair('a', "aaaafoaobAAaAAArAssAa"))
}

func killPair(toKill rune, line string) string {
	res := make([]string, 0)
	last := 0
	for i, c := range(strings.ToLower(line)) {
		if c == toKill {
			res = append(res, line[last:i])
			last = i+1
		}
	}
	res = append(res, line[last:])
	return strings.Join(res, "")
}
