package night

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type LogLine struct {
	Spec  string
	Date  string
	Hours int
	Mins  int
	Ev    Event
}

func LogGenerator(lines []string) chan GuardState {
	c := make(chan GuardState)
	go func() {
		var currentGuard int
		var sleeping bool
		var t int
		for _, s := range lines {
			l, err := NewLogLine(s)
			if err != nil {
				panic(err)
			}
			// c <- l // temp
			if l.Ev.Sleep() {
				for ; t < l.Mins; t++ {
					c <- GuardState{guard: currentGuard, asleep: sleeping, time: t}
				}
				assert(currentGuard != 0)
				assert(!sleeping)
				sleeping = true
			} else if l.Ev.Wake() {
				for ; t < l.Mins; t++ {
					c <- GuardState{guard: currentGuard, asleep: sleeping, time: t}
				}
				assert(currentGuard != 0)
				assert(sleeping)
				sleeping = false
			} else if nextGuard := l.Ev.GuardChange(); nextGuard != 0 {
				if currentGuard != 0 {
					for ; t < 60; t++ {
						c <- GuardState{guard: currentGuard, asleep: sleeping, time: t}
					}
				}
				assert(!sleeping) // TODO how early/late can they arrive?
				currentGuard = nextGuard
				sleeping = false // redundant
				if l.Hours > 0 {
					t = 0
				} else {
					t = l.Mins
				}
			} else {
				panic(fmt.Sprintf("bad event: %#v", l.Ev))
			}
		}
		close(c)
	}()
	return c
}

type GuardState struct {
	guard  int
	asleep bool
	time   int
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func NewLogLine(spec string) (res LogLine, err error) {
	c := stringChomper{src: spec, pos: 0}
	c.expect("[")
	date := c.chomp(len("YYYY-MM-DD"))
	c.expect(" ")
	_hours := c.chomp(len("hh"))
	c.expect(":")
	_mins := c.chomp(len("mm"))
	c.expect("] ")
	_event := c.chompRest()

	hours, err := strconv.Atoi(_hours)
	if err != nil {
		return
	}
	mins, err := strconv.Atoi(_mins)
	if err != nil {
		return
	}
	event, err := newEvent(_event)
	if err != nil {
		return
	}

	res = LogLine{
		Spec:  spec,
		Date:  date,
		Mins:  mins,
		Hours: hours,
		Ev:    event,
	}
	return
}

type stringChomper struct {
	src string
	pos int
}

func (c *stringChomper) expect(e string) {
	if actual := c.chomp(len(e)); actual != e {
		panic(fmt.Sprintf("Expected to chomp '%s', got '%s'", e, actual))
	}
}

func (c *stringChomper) chomp(n int) string {
	newPos := c.pos + n
	assert(newPos <= len(c.src))
	res := c.src[c.pos:newPos]
	c.pos = newPos
	return res
}

func (c *stringChomper) chompRest() string {
	return c.chomp(len(c.src) - c.pos)
}

// union; implementors must return exactly one non-"zero"
//   value from the three interface methods
type Event interface {
	GuardChange() int
	Wake() bool
	Sleep() bool
}

type defaultEvent struct{}

func (e defaultEvent) GuardChange() int {
	return 0
}
func (e defaultEvent) Wake() bool {
	return false
}
func (e defaultEvent) Sleep() bool {
	return false
}

type guardChangeEvent struct {
	defaultEvent
	n int
}

func (e guardChangeEvent) GuardChange() int {
	return e.n
}

type wakeEvent struct {
	defaultEvent
}

func (e wakeEvent) Wake() bool {
	return true
}

type sleepEvent struct {
	defaultEvent
}

func (e sleepEvent) Sleep() bool {
	return true
}

func newEvent(s string) (res Event, err error) {
	switch {
	case s == "wakes up":
		res = wakeEvent{}
	case s == "falls asleep":
		res = sleepEvent{}
	default:
		if n := guardChangeId(s); n != 0 {
			res = guardChangeEvent{n: n}
		}
	}
	if res == nil {
		err = errors.New(fmt.Sprintf("unrecognized event: '%s'", s))
	}
	return
}

// panics
// 0 = no match
// other = guard id found
func guardChangeId(s string) int {
	re := regexp.MustCompile("^Guard #(\\d+) begins shift$")
	groups := re.FindStringSubmatch(s)
	if groups == nil {
		return 0
	}
	num, err := strconv.Atoi(groups[1])
	if err != nil {
		panic(err)
	}
	if num == 0 {
		panic(fmt.Sprintf("Guard cannot be id 0: '%s'", s))
	}

	return num
}
