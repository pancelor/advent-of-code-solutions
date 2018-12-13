package main

import (
	"bufio"
	debug "log"
	"io"
	"fmt"
	"os"
	"sort"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

	answer, err := solve(os.Stdin)
	check(err)
	fmt.Println(answer)
}

func solve(in io.Reader) (answer string, err error) {
	data := make([]req, 0)
	seen := make(map[byte]bool)
	{
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			line := scanner.Text()
			req := req{}
			fmt.Sscanf(line, "Step %c must be finished before step %c can begin.", &req.a, &req.b)
			data = append(data, req)
			seen[req.a] = true
			seen[req.b] = true
		}
		err = scanner.Err()
		if err != nil {
			return
		}

		// for _, v := range data {
		// 	debug.Println(v)
		// }
	}

	ready := make(map[byte]bool)
	postreqs := make(connections)
	prereqs := make(connections)
	{
		for k := range seen {
			ready[k] = true
			postreqs[k] = make(map[byte]bool)
			prereqs[k] = make(map[byte]bool)
		}

		for _, v := range data {
			ready[v.b] = false
			postreqs[v.a][v.b] = true // once set to true this will never be set false
			prereqs[v.b][v.a] = true // these _will_ be set back to false as prerequisites complete
		}
	}

	// debug.Println("ready")
	// debugSet(ready)
	// debug.Println("postreqs")
	// debugConnections(postreqs)
	// debug.Println("prereqs")
	// debugConnections(prereqs)

	var s strings.Builder
	{
		for {
			ord := arrFromSet(ready)
			sort.Sort(sortBytes(ord))
			if len(ord) == 0 {
				break
			}
			r := ord[0]

			fmt.Fprintf(&s, "%c", r)
			ready[r] = false
			// debug.Printf("postreqs[%c]\n", r)
			// debugSet(postreqs[r])
			for k := range postreqs[r] {
				// debug.Printf("analyzing %c\n", k)
				prereqs[k][r] = false
				remaining := arrFromSet(prereqs[k])
				// debug.Printf("remaining:")
				// debugSet(prereqs[k])
				if len(remaining) == 0 {
					ready[k] = true
				}
			}
		}
	}

	answer = s.String()
	return
}

type connections map[byte]map[byte]bool

type sortBytes []byte
func (s sortBytes) Less(i, j int) bool {
	return s[i] < s[j]
}
func (s sortBytes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortBytes) Len() int {
	return len(s)
}

func arrFromSet(c map[byte]bool) []byte{
	arr := make([]byte, 0)
	for k, v := range c {
		if v {
			arr = append(arr, k)
		}
	}
	return arr
}

func debugSet(c map[byte]bool) {
	// fmt.Printf("ready{")
	fmt.Printf("{")
	for k, v := range c {
		if v {
			fmt.Printf("%c", k)
		}
	}
	fmt.Printf("}\n")
	// fmt.Printf("not-ready{")
	// for k, v := range c {
	// 	if !v {
	// 		fmt.Printf("%c", k)
	// 	}
	// }
	// fmt.Printf("}\n")
}

func debugConnections(c connections) {
	fmt.Printf("{\n")
	for k := range c {
		if len(c[k]) == 0 {
			continue
		}
		fmt.Printf("\t%c->{", k)
		for k2 := range c[k] {
			fmt.Printf("%c", k2)
			if v, ok := c[k][k2]; v == false && ok {
				fmt.Printf("(met)")
			}
		}
		fmt.Printf("}\n")
	}
	fmt.Printf("}\n")
}

// "<a> must be finished before <b> can begin."
type req struct{
	a, b byte
}

func (p req) String() string {
	return fmt.Sprintf("%c->%c", p.a, p.b)
}
