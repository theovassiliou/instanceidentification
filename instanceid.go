package instanceid

import (
	"strconv"
	"strings"
	"time"
)

// This package supports the handling of instance-identification fields
// as proposed in Theos paper

// An instance is represented through it MIID
// MIID := <sN> "/" <vN> ["/" <vA>] "%" <t>s
// Example:
//		msA/1.1/feature-branch-2345abcd%222s
// The complete call-graph including it's own MIID
// is represented by:
// CIID := MIID [ "(" UIDs+ ")"]
// UIDs := CIID [ "+" CIID ]+

// This package provides some helpers to work with this
// type of instance-identification

// CIID := MIID [ "(" UIDs+ ")"]
// UIDs := CIID [ "+" CIID ]+
// MIID := <sN> "/" <vN> ["/" <vA>] "%" <t>s

// Ciid represents the complete call-graph as instance-id
type Ciid struct {
	Miid  Miid
	Ciids Stack
}

// Miid represents the instance only by it's name, version, additional information
// and epoch time
type Miid struct {
	Sn string
	Vn string
	Va string
	T  int
}

// Stack represents a list of services that have been called by the Ciid
type Stack []Ciid

// NewCiid creates a new Ciid from a string in the form of
// Sn1/Vn1/Va1%t1s(Sn2/Vn2/Va2%t2s+Sn3/Vn3/Va3%t3s(Sn4/Vn4/Va4%t4s))
func NewCiid(id string) (ciid Ciid) {
	return parseCiid(id)
}

// NewMiid creates a new Miid from a string in the of
// Sn1/Vn1/Va1%t1s
// in case a Ciid is being provided the Miid part is only
// returned
// If there are syntax errors an empty Miid will be returned
func NewMiid(id string) (miid Miid) {
	return parseMIID(id)
}

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
		return Ciid{}, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

// SanityCheck checks the given miid against some rules to ensure that it can be an Miid
// returns true if miid could be an Miid false otherwise
func SanityCheck(miid string) bool {
	miid = strings.TrimSpace(miid)

	// last rune must be s
	if !strings.HasSuffix(miid, "s") {
		return false
	}

	// at least two /
	s := strings.Count(miid, "/")
	if s < 1 || s > 3 {
		return false
	}

	// no '+', no '(' no ')'

	if strings.ContainsAny(miid, "+()") {
		return false
	}

	return true
}

func (myself *Ciid) SetStack(callStack Stack) *Ciid {
	myself.Ciids = callStack
	return myself
}

func (myself *Ciid) ClearStack() *Ciid {
	myself.Ciids = nil
	return myself
}

func (c *Ciid) String() string {
	sB := strings.Builder{}
	sB.WriteString(c.Miid.String())
	if len(c.Ciids) > 0 {
		sB.WriteString("(")
		for i, a := range c.Ciids {
			sB.WriteString(a.String())
			if i+1 < len(c.Ciids) {
				sB.WriteString("+")
			}
		}
		sB.WriteString(")")
	}

	return sB.String()
}

func (m *Miid) String() string {
	sB := strings.Builder{}
	if m.Sn != "" {
		sB.WriteString(m.Sn)
		if m.Vn != "" {
			sB.WriteString("/" + m.Vn)
		}
		if m.Va != "" {
			sB.WriteString("/" + m.Va)
		}
		sB.WriteString("%" + strconv.Itoa(m.T) + "s")
	}
	return sB.String()
}

// Contains returns true if the Ciid contains the left aligned miid as part of the call graph
func (ciid *Ciid) Contains(miid string) bool {
	if miid == "" {
		return false
	}
	return strings.Contains(ciid.String(), miid)
}

func (ciid *Ciid) SetEpoch(startTime time.Time) *Ciid {
	epoch := time.Since(startTime)
	ciid.Miid.T = int(epoch.Seconds())
	return ciid
}

func (miid *Miid) SetEpoch(startTime time.Time) *Miid {
	epoch := time.Since(startTime)
	miid.T = int(epoch.Seconds())
	return miid
}

// Contains returns true if s is contained left aligned, else or if s is empty return false
func (m *Miid) Contains(s string) bool {
	if s == "" {
		return false
	}
	return strings.Contains(m.String(), s)
}

func parseCiid(id string) (ciid Ciid) {
	name, arg := seperateFNameFromArg(id)

	if arg == "" {
		return Ciid{Miid: parseMIID(name)}
	}
	me := Ciid{Miid: parseMIID(name)}
	me.Ciids = parseArguments(arg)
	return me
}

func parseArguments(arg string) (ciids []Ciid) {
	ss := splitOnPlus(arg)

	for _, a := range ss {
		ciids = append(ciids, parseCiid(a))
	}
	return ciids
}

func parseMIID(id string) (miid Miid) {
	if !SanityCheck(id) {
		return miid
	}
	s := strings.SplitN(id, "/", -1)
	l := len(s)
	var r Miid
	if l >= 1 {
		r.Sn = s[0]
	}
	if l == 2 {
		e := strings.Split(s[1], "%")
		r.Vn = e[0]
		if len(e) > 1 {
			t, _ := strconv.Atoi(strings.Split(e[1], "s")[0])
			r.T = t
		}
	} else if l >= 2 {
		r.Vn = s[1]
	}

	if l == 3 {
		e := strings.Split(s[2], "%")
		r.Va = e[0]
		t, _ := strconv.Atoi(strings.Split(e[1], "s")[0])
		r.T = t
	}

	return r
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

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
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
