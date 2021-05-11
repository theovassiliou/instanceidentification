package instanceid

import (
	"strings"
	"time"
)

type StdCiid struct {
	miid  Miid
	ciids Stack
}

// NewCiid creates a new Ciid from a string in the form of
// Sn1/Vn1/Va1%t1s(Sn2/Vn2/Va2%t2s+Sn3/Vn3/Va3%t3s(Sn4/Vn4/Va4%t4s))
func NewStdCiid(id string) (ciid *StdCiid) {
	return parseCiid(id)
}

func (c StdCiid) Miid() Miid {
	return c.miid
}

func (c StdCiid) Ciids() Stack {
	return c.ciids
}

func (c *StdCiid) SetCiids(s Stack) Ciid {
	c.ciids = s
	return c
}

// SetEpoch sets the epoch field based on a given StartTime. Chainable.
func (ciid *StdCiid) SetEpoch(startTime time.Time) Ciid {
	epoch := time.Since(startTime)
	ciid.miid.SetT(int(epoch.Seconds()))
	return ciid
}

// String returns the textual representation of the Ciid
func (c StdCiid) String() string {
	sB := strings.Builder{}
	sB.WriteString(c.miid.String())
	if len(c.ciids) > 0 {
		sB.WriteString("(")
		for i, a := range c.ciids {
			sB.WriteString(a.String())
			if i+1 < len(c.ciids) {
				sB.WriteString("+")
			}
		}
		sB.WriteString(")")
	}

	return sB.String()
}

// Contains returns true if the Ciid contains the left aligned miid as part of the call graph
func (c StdCiid) Contains(miid string) bool {
	if miid == "" {
		return false
	}
	return strings.Contains(c.String(), miid)
}

func parseCiid(id string) *StdCiid {
	me := new(StdCiid)
	name, arg := seperateFNameFromArg(id)

	me.miid = parseMIID(name)

	if arg == "" || strings.Contains(arg, "+") {
		return me
	}

	me.ciids = parseArguments(arg)
	return me
}

func seperateFNameFromArg(signature string) (name, arg string) {
	n := strings.Builder{}
	a := strings.Builder{}
	var inArgs bool
	var count int
	for _, s := range signature {
		if s == '(' {
			count++
			inArgs = true
		}
		if s == ')' {
			count--
		}

		if !inArgs {
			n.WriteRune(s)
		} else if count == 1 && s != '(' {
			a.WriteRune(s)
		} else if count > 1 {
			a.WriteRune(s)
		}
	}

	return n.String(), a.String()
}

func parseArguments(arg string) (ciids Stack) {
	ss := splitOnPlus(arg)

	for _, a := range ss {
		ciids = append(ciids, parseCiid(a))
	}
	return ciids
}

func splitOnPlus(s string) (ss []string) {
	var openClose int
	var splitPos []int

	// first find parenthis pairs
	for pos, char := range s {
		if char == '(' {
			openClose++
		} else if char == ')' {
			openClose--
		} else if char == '+' {
			if openClose == 0 {
				splitPos = append(splitPos, pos)
			}
		}
	}

	//   split arguments
	if len(splitPos) > 0 {
		oldPos := -1
		for _, s2 := range splitPos {
			ss = append(ss, s[oldPos+1:s2])
			oldPos = s2
		}
		ss = append(ss, s[oldPos+1:])
	} else {
		ss = append(ss, s)
	}
	return deleteEmpty(ss)
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// -- Stack Implementation for StdCiid

// Push a new value onto the stack
func (s *Stack) Push(str Ciid) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (Ciid, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

func (m StdCiid) SetStack(callStack Stack) StdCiid {
	m.SetCiids(callStack)
	return m
}

func (m StdCiid) ClearStack() StdCiid {
	m.SetCiids(nil)
	return m
}
