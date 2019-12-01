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
	total := 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		v := scanner.Text()

		mass, err := strconv.Atoi(v)
		if err != nil {
			return "", err
		}
		total += fuelForMassSimple(mass)
	}
	err := scanner.Err()
	if err != nil {
		return "", err
	}

	total += fuelForMass(total) // fuel itself has mass

	return strconv.Itoa(total), nil
}

func fuelForMass(mass int) int {
	total := 0
	for {
		fuel := fuelForMassSimple(mass)
		if fuel <= 0 {
			break
		}
		total += fuel
		mass = fuel
	}
	return total
}

func fuelForMassSimple(mass int) int {
	return (mass / 3) - 2
}
