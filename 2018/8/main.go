// part 1 is done; haven't tried submitting yet

package main

import (
	"fmt"
	"io"
	"os"
)

// const instring = "2 3 0 3 10 11 12 1 1 0 1 99 2 1 1 2"

func main() {
	answer, err := solve(os.Stdin)
	// answer, err := solve(strings.NewReader(instring))
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

func readInt(in io.Reader) (res int, err error) {
	_, err = fmt.Fscanf(in, "%d", &res)
	return
}

func solve(in io.Reader) (answer int, err error) {
	root, err := readNode(in)
	if err != nil {
		return
	}
	// log.Printf("%#v\n", root)
	answer = reduceNode(func(reducedChildren []int, metadata []int) int {
		return sum(reducedChildren) + sum(metadata)
	}, root)
	return
}

func sum(a []int) int {
	res := 0
	for _, x := range a {
		res += x
	}
	return res
}

// fn takes 1: reduced children and 2: metadata
func reduceNode(fn func([]int, []int) int, n Node) int {
	results := make([]int, len(n.children))
	for _, c := range n.children {
		results = append(results, reduceNode(fn, c))
	}
	return fn(results, n.metadata)
}

func readNode(in io.Reader) (node Node, err error) {
	nChildren, err := readInt(in)
	if err != nil {
		return
	}
	nMetadata, err := readInt(in)
	if err != nil {
		return
	}
	node = NewNode(nChildren, nMetadata)
	for i := 0; i < nChildren; i++ {
		child, err2 := readNode(in)
		if err2 != nil {
			err = err2
			return
		}
		node.children[i] = child
	}
	for i := 0; i < nMetadata; i++ {
		meta, err2 := readInt(in)
		if err2 != nil {
			err = err2
			return
		}
		node.metadata[i] = meta
	}
	return
}

type Node struct {
	children []Node
	metadata []int
}

func NewNode(nChildren, nMetadata int) Node {
	return Node{
		children: make([]Node, nChildren),
		metadata: make([]int, nMetadata),
	}
}
