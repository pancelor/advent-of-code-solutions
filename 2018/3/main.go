package main

import (
	"bufio"
	debug "log"
	"io"
	"fmt"
	"os"
	"./claim"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	debug.SetFlags(0)
	debug.SetPrefix("debug: ")

  answer, err := solve(newValReader(os.Stdin))
  if err != nil {
    panic(err)
  }
  fmt.Println(answer)
}

func solve(reader valReader) (answer int, outErr error) {
	claims := make([]claim.Claim, 0)
	for v := range reader.Vals() {
		c := claim.New(v)
		// x1, y1 := c.TopLeft()
		// x2, y2 := c.BottomRightExclusive()
		// debug.Println(x1, y1, x2, y2)
		claims = append(claims, c)
	}
	outErr = reader.Err()
	if outErr != nil {
		return
	}

	w, h := getMaxSize(claims)

	fmt.Println(claims)
	fmt.Println(w, h)

	return
}

func getMaxSize(claims []claim.Claim) (int, int) {
	return 1,2
}

type projection func(c claim.Claim) int
func maxByKey(arr []claim.Claim, key projection, int default) claim.Claim {
	res = 0
	for _, e := range arr {
		val := key(e)
		if val > res {
			res = val
		}
	}
	return res
}


// valReader converts an `io.Reader` to a `chan string`
// usage:
//   reader := newValReader(os.Stdin)
//   for val := reader.Vals() {
//     _ = val
//   }
//   if reader.Err() != nil {
//     panic(err)
//   }
type valReader struct {
	in io.Reader
	out chan string
	err error
}

func newValReader(in io.Reader) valReader {
	out := make(chan string)
	return valReader{
		in: in,
		out: out,
	}
}

func (r *valReader) Vals() chan string {
	go func() {
		scanner := bufio.NewScanner(r.in)
		for scanner.Scan() {
			r.out <- scanner.Text()
		}
		if r.check(scanner.Err()) {
			return
		}
		close(r.out)
	}()
	return r.out
}

func (r *valReader) Err() error {
	return r.err
}

func (r *valReader) check(err error) bool {
	if err != nil {
		r.err = err
		close(r.out)
		return true
	}
	return false
}
