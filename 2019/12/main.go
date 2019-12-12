package main

import (
	"fmt"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

type point struct {
	x, y, z int
}

func (p point) add(o point) point {
	return point{
		x: p.x + o.x,
		y: p.y + o.y,
		z: p.z + o.z,
	}
}

func (p point) sub(o point) point {
	return point{
		x: p.x - o.x,
		y: p.y - o.y,
		z: p.z - o.z,
	}
}

func (p point) sign() point {
	res := point{}
	if p.x < 0 {
		res.x = -1
	} else if p.x > 0 {
		res.x = 1
	}
	if p.y < 0 {
		res.y = -1
	} else if p.y > 0 {
		res.y = 1
	}
	if p.z < 0 {
		res.z = -1
	} else if p.z > 0 {
		res.z = 1
	}

	return res
}

func printState(moons []point, vels []point) {
	for i := 0; i < len(moons); i++ {
		fmt.Printf("pos=%v vel=%v\n", moons[i], vels[i])
	}
	fmt.Printf("\n")
}

func applyGravity(moons []point, vels []point) {
	for i := 0; i < len(moons); i++ {
		for j := 0; j < len(moons); j++ {
			if i == j {
				continue
			}
			sign := moons[j].sub(moons[i]).sign()
			vels[i] = vels[i].add(sign)
		}
	}
}

func applyVelocity(moons []point, vels []point) {
	for i := 0; i < len(moons); i++ {
		moons[i] = moons[i].add(vels[i])
	}
}

func energy(moons []point, vels []point) int {
	total := 0
	for i := 0; i < len(moons); i++ {
		pot := 0
		pot += helpers.Abs(moons[i].x)
		pot += helpers.Abs(moons[i].y)
		pot += helpers.Abs(moons[i].z)
		kin := 0
		kin += helpers.Abs(vels[i].x)
		kin += helpers.Abs(vels[i].y)
		kin += helpers.Abs(vels[i].z)
		// fmt.Println(moons[i], pot, vels[i], kin)
		total += pot * kin
	}
	return total
}

func solve(moons []point) interface{} {
	vels := []point{
		point{}, point{}, point{}, point{},
	}
	assert(len(moons) == len(vels), "mismatched len")

	for i := 0; i < 1000; i++ {
		// if i == 0 {
		// 	fmt.Printf("After %d steps:\n", i)
		// 	printState(moons, vels)
		// }
		applyGravity(moons, vels)
		applyVelocity(moons, vels)
		// fmt.Printf("After %d steps:\n", i+1)
		// printState(moons, vels)
	}

	return energy(moons, vels)
}

func init() {
	// tests go here
	assert(true, "t1")

	// assert(false, "exit after tests")
}

func main() {
	input, err := getInput()
	check(err)
	answer := solve(input)
	fmt.Printf("answer:\n%v\n", answer)
}

func getInput() ([]point, error) {
	return []point{
		// ex 1
		// point{x: -1, y: 0, z: 2},
		// point{x: 2, y: -10, z: -7},
		// point{x: 4, y: -8, z: 8},
		// point{x: 3, y: 5, z: -1},
		// ex 2
		// point{x: -8, y: -10, z: 0},
		// point{x: 5, y: 5, z: 10},
		// point{x: 2, y: -7, z: 3},
		// point{x: 9, y: -8, z: -3},
		// real input
		point{x: 0, y: 6, z: 1},
		point{x: 4, y: 4, z: 19},
		point{x: -11, y: 1, z: 8},
		point{x: 2, y: 19, z: 15},
	}, nil
}
