package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"log"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("debug: ")

	scanner := bufio.NewScanner(os.Stdin)
	var total int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		// log.Println(val)
		check(err)
		total += val
	}
	fmt.Println(total)
}
