package main

import (
	"fmt"
	"strings"

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

func (p point) multHadamard(o point) point {
	return point{
		x: p.x * o.x,
		y: p.y * o.y,
		z: p.z * o.z,
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

func (state *simulationState) String() string {
	var s strings.Builder
	for i := 0; i < len(state.moons); i++ {
		fmt.Fprintf(&s, "pos=%v vel=%v\n", state.moons[i], state.vels[i])
	}
	fmt.Fprintf(&s, "\n")
	return s.String()
}

func (state *simulationState) applyGravity() {
	for i := 0; i < len(state.moons); i++ {
		for j := 0; j < len(state.moons); j++ {
			if i == j {
				continue
			}
			sign := state.moons[j].sub(state.moons[i]).sign()
			state.vels[i] = state.vels[i].add(sign)
		}
	}
}

func (state *simulationState) applyVelocity() {
	for i := 0; i < len(state.moons); i++ {
		state.moons[i] = state.moons[i].add(state.vels[i])
	}
}

func (state *simulationState) energy() int {
	total := 0
	for i := 0; i < len(state.moons); i++ {
		pot := 0
		pot += helpers.Abs(state.moons[i].x)
		pot += helpers.Abs(state.moons[i].y)
		pot += helpers.Abs(state.moons[i].z)
		kin := 0
		kin += helpers.Abs(state.vels[i].x)
		kin += helpers.Abs(state.vels[i].y)
		kin += helpers.Abs(state.vels[i].z)
		// fmt.Println(state.moons[i], pot, state.vels[i], kin)
		total += pot * kin
	}
	return total
}

// func hash(state) []byte {
// 	first := sha256.New()
// 	return first.Sum(nil)
// }

type simulationState struct {
	moons [4]point
	vels  [4]point
}

func solve(realState simulationState) interface{} {
	var stateX, stateY, stateZ simulationState
	for i := 0; i < len(realState.moons); i++ {
		stateX.moons[i] = realState.moons[i].multHadamard(point{1, 0, 0})
		stateY.moons[i] = realState.moons[i].multHadamard(point{0, 1, 0})
		stateZ.moons[i] = realState.moons[i].multHadamard(point{0, 0, 1})
		stateX.vels[i] = realState.vels[i].multHadamard(point{1, 0, 0})
		stateY.vels[i] = realState.vels[i].multHadamard(point{0, 1, 0})
		stateZ.vels[i] = realState.vels[i].multHadamard(point{0, 0, 1})
	}
	x, ok := solveSlow(stateX, 1000000)
	assert(ok, "x bound too low")
	y, ok := solveSlow(stateY, 1000000)
	assert(ok, "y bound too low")
	z, ok := solveSlow(stateZ, 1000000)
	assert(ok, "z bound too low")
	fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)
	return helpers.Lcm(x, y, z)
}

func solveSlow(state simulationState, limit int) (int, bool) {
	seen := make(map[simulationState]bool)
	for i := 0; i < limit; i++ {
		// if i == 0 {
		// 	fmt.Printf("After %d steps:\n%s", i, state.String())
		// }

		if seen[state] {
			return i, true
		}
		seen[state] = true
		state.applyGravity()
		state.applyVelocity()
		// fmt.Printf("After %d steps:\n%s", i, state.String())
	}

	return 0, false
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

func getInput() (simulationState, error) {
	moons := [4]point{
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
	}
	s := simulationState{moons: moons}
	return s, nil
}
