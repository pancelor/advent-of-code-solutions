package claim

import (
	"fmt"
)

type Claim struct {
	spec string
	id int
	x int
	y int
	w int
	h int
}

// spec example: "#2 @ 391,45: 27x20"
func New(spec string) Claim {
	c := Claim{
		spec: spec,
	}
	fmt.Sscanf(spec, "#%d @ %d,%d: %dx%d", &c.id, &c.x, &c.y, &c.w, &c.h)
	if c.w == 0 || c.h == 0 {
		panic(fmt.Sprintf("zero area claim %v", spec))
	}
	return c
}

func (c Claim) Id() int {
	return c.id
}

func (c Claim) X() int {
	return c.x
}

func (c Claim) Y() int {
	return c.y
}

func (c Claim) W() int {
	return c.w
}

func (c Claim) H() int {
	return c.h
}

// one square past, actually
func (c Claim) Bottom() int {
	return c.y + c.h
}

// one square past, actually
func (c Claim) Right() int {
	return c.x + c.w
}
