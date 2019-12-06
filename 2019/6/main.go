package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// KidMap .
type KidMap map[string]map[string]bool

// ParentMap .
type ParentMap map[string]string

func buildKidMap(input []Orbit) KidMap {
	kidMap := make(map[string]map[string]bool)
	for _, orbit := range input {
		// fmt.Printf("orbit=%#v\n", orbit)
		_, ok := kidMap[orbit.center]
		if !ok {
			kidMap[orbit.center] = make(map[string]bool)
		}
		kidMap[orbit.center][orbit.other] = true
	}
	return kidMap
}

func buildParentMap(input []Orbit) ParentMap {
	parentMap := make(map[string]string)
	for _, orbit := range input {
		// fmt.Printf("orbit=%#v\n", orbit)
		parentMap[orbit.other] = orbit.center
	}
	return parentMap
}

func solve1(input []Orbit) (interface{}, error) {
	kidMap := buildKidMap(input)

	total := 0
	nextLevel := []string{"COM"}
	for depth := 0; len(nextLevel) > 0; depth++ {
		fmt.Printf("nextLevel=%#v\n", nextLevel)
		currLevel := nextLevel
		nextLevel = []string{}
		for _, s := range currLevel {
			total += depth
			for k := range kidMap[s] {
				nextLevel = append(nextLevel, k)
			}
		}
	}

	return total, nil
}

func solve2(input []Orbit) (interface{}, error) {
	parentMap := buildParentMap(input)

	santaDist := make(map[string]int)

	{ // find YOU chain
		node := "SAN"
		for i := 0; true; i++ {
			// fmt.Printf("SAN; node=%#v\n", node)
			var ok bool
			node, ok = parentMap[node]
			if ok == false {
				break
			}
			santaDist[node] = i
		}
	}

	{ // find where YOU chain connects to SAN chain
		node := "YOU"
		for i := 0; true; i++ {
			// fmt.Printf("YOU; node=%#v\n", node)
			var ok bool
			node, ok = parentMap[node]
			if !ok {
				return nil, fmt.Errorf("Didn't find santa :(")
			}
			dist, ok := santaDist[node]
			if ok {
				return dist + i, nil
			}
		}
	}

	return nil, fmt.Errorf("Didn't find santa... and also somehow escaped the loop??")
}

func test() {
	assert(true, "t1")

	// assert(false, "exit after tests")
}

func main() {
	test()

	input, err := getInput()
	if err != nil {
		panic(err)
	}

	// answer, err := solve1(input)
	answer, err := solve2(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("answer:\n%v\n", answer)
}

// Orbit .
type Orbit struct {
	center string
	other  string
}

//
// helpers
//

// for temporary use only
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// for temporary use only
func assert(b bool, msg string) {
	if !b {
		panic(errors.New(msg))
	}
}

var source = os.Stdin

func getInput() ([]Orbit, error) {
	lines, err := getLines()
	if err != nil {
		return nil, err
	}

	var res []Orbit
	for _, l := range lines {
		assert(len(l) == 7, "bad input")
		res = append(res, Orbit{center: l[0:3], other: l[4:7]})
	}

	return res, nil
}

func getLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(source)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
