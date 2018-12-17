package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
	fmt.Printf("%s\n", root)
	answer = reduceNode(func(reducedChildren []int, metadata []int) (res int) {
		// if len(reducedChildren) == 0 {
		// 	res = sum(metadata)
		// } else {
		// 	for _, m := range metadata {
		// 		res += reducedChildren[m-1]
		// 	}
		// }
		res = sum(reducedChildren) + sum(metadata)
		log.Printf("reduce(%v, %v)=%v", reducedChildren, metadata, res)
		return
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
	for i, c := range n.children {
		results[i] = reduceNode(fn, c)
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

var indent int

type indentPrinter struct {
	f io.Writer
}

func newIndentPrinter(f io.Writer) indentPrinter {
	return indentPrinter{f: f}
}

func (p *indentPrinter) Printf(s string, a ...interface{}) {
	fmt.Fprintf(p.f, s, a...)
}

func (p *indentPrinter) Println(s ...interface{}) {
	fmt.Fprintln(p.f, s...)
	fmt.Fprintf(p.f, strings.Repeat("\t", indent))
}

func (n Node) String() string {
	var s strings.Builder
	p := newIndentPrinter(io.Writer(&s))
	p.Printf("Node{{")
	if len(n.children) > 0 {
		indent += 1
		for _, c := range n.children {
			p.Println()
			p.Printf("%s", c)
		}
		indent -= 1
		p.Println()
	}
	p.Printf("}[")
	for i, m := range n.metadata {
		if i != 0 {
			p.Printf(" ")
		}
		p.Printf("%d", m)
	}
	p.Printf("]}")
	return s.String()
}
