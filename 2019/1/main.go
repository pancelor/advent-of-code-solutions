package main

import (
	"bufio"
	"fmt"
	"io"
	debug "log"
	"os"
	"strconv"
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

func solve(reader io.Reader) (string, error) {
	counter := 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		v := scanner.Text()

		mass, err := strconv.Atoi(v)
		if err != nil {
			return "", err
		}
		counter += fuelForMass(mass)
	}
	err := scanner.Err()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(counter), nil
}

func fuelForMass(mass int) int {
	return (mass / 3) - 2
}
