package main

import (
	"fmt"
)

func main() {
	g :=foo()
	for i := range g {
		fmt.Println(i)
	}
}

func foo() chan int {
	c := make(chan int)
	go func() {
		c<-1
		c<-2
		c<-3
		close(c)
	}()
	return c
}