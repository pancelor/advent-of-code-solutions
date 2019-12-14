package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/pancelor/advent-of-code-solutions/2019/helpers"
)

var assert = helpers.Assert
var check = helpers.Check

// Input is the type of input we pass to solve()
type Input []Reaction

type Reaction struct {
	lhs []ChemicalAmount
	rhs ChemicalAmount
}

type ChemicalAmount struct {
	n        int
	chemical Chemical
}

type Chemical string

type InvertedReaction struct {
	n   int
	lhs []ChemicalAmount
}

type Balance map[Chemical]int

func makeBalance() Balance {
	return make(map[Chemical]int)
}

// second arg is whether its done or not
func (b Balance) nextRequirement() (ChemicalAmount, bool) {
	for k, v := range b {
		if v < 0 && k != Chemical("ORE") {
			return ChemicalAmount{chemical: k, n: v}, false
		}
	}
	return ChemicalAmount{}, true
}

func (b Balance) need(ca ChemicalAmount, m int) {
	b[ca.chemical] -= ca.n * m
}

func (b Balance) make(ca ChemicalAmount, m int) {
	b[ca.chemical] += ca.n * m
}

func solve(in Input) interface{} {
	book := make(map[Chemical]InvertedReaction)
	for _, rxn := range in {
		_, prs := book[rxn.rhs.chemical]
		assert(!prs, "multiply defined output %s", rxn.rhs.chemical)
		book[rxn.rhs.chemical] = InvertedReaction{n: rxn.rhs.n, lhs: rxn.lhs}
	}

	// fmt.Printf("val=%#v\n", book)

	bal := makeBalance()
	bal.need(ChemicalAmount{chemical: Chemical("FUEL"), n: 1}, 1)
	for {
		next, done := bal.nextRequirement()
		fmt.Printf("Making %#v (%v)\n", next, done)
		if done {
			break
		}
		ir := book[next.chemical]
		assert(ir.n != 0, "don't know how to make %s", next.chemical)
		fmt.Printf("ir=%#v\n", ir)
		times := int(math.Ceil(math.Abs(float64(next.n) / float64(ir.n))))
		fmt.Println(times, next.n, ir.n)
		bal.make(ChemicalAmount{chemical: next.chemical, n: ir.n}, times)
		for _, ca := range ir.lhs {
			bal.need(ca, times)
		}
	}

	fmt.Println(bal)
	return bal[Chemical("ORE")]
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

func getInput() (Input, error) {
	lines, err := helpers.GetLines()
	if err != nil {
		return nil, err
	}

	var res []Reaction
	for _, line := range lines {
		if line == "" {
			continue
		}

		sides := strings.Split(line, "=>")
		assert(len(sides) == 2, "bad num sides")
		lhsStr := strings.Split(sides[0], ",")
		rhsStr := sides[1]
		var lhs []ChemicalAmount
		for _, in := range lhsStr {
			lhs = append(lhs, parseChemicalAmount(in))
		}
		rhs := parseChemicalAmount(rhsStr)

		res = append(res, Reaction{lhs: lhs, rhs: rhs})
	}

	return res, nil
}

func parseChemicalAmount(s string) ChemicalAmount {
	var n int
	var chemical string
	fmt.Sscanf(strings.Trim(s, " "), "%d %s", &n, &chemical)
	return ChemicalAmount{
		chemical: Chemical(chemical),
		n:        n,
	}
}
