package instanceid

import (
	"strings"
	"time"
)

// Ciid represents the complete call-graph as instance-id
type Ciid struct {
	Miid   Miid
	Ciids  Stack
	decode DecodingFunction
}

// NewCiid creates a new Ciid from a string in the form of
// Sn1/Vn1/Va1%t1s(Sn2/Vn2/Va2%t2s+Sn3/Vn3/Va3%t3s(Sn4/Vn4/Va4%t4s))
func NewCiid(id string) (ciid Ciid) {
	return parseCiid(id)
}

func (c *Ciid) WithDecoding(d DecodingFunction) *Ciid {
	c.decode = d
	c.Miid.WithDecoding(d)
	for i := range c.Ciids {
		c.Ciids[i].WithDecoding(d)
	}
	return c
}

// String returns the textual representation of the Ciid
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

// Contains returns true if the Ciid contains the left aligned miid as part of the call graph
func (ciid *Ciid) Contains(miid string) bool {
	if miid == "" {
		return false
	}
	return strings.Contains(ciid.String(), miid)
}

// SetEpoch sets the epoch field based on a given StartTime. Chainable.
func (ciid *Ciid) SetEpoch(startTime time.Time) *Ciid {
	epoch := time.Since(startTime)
	ciid.Miid.T = int(epoch.Seconds())
	return ciid
}

func parseCiid(id string) (ciid Ciid) {
	name, arg := seperateFNameFromArg(id)

	me := Ciid{Miid: parseMIID(name)}

	if arg == "" {
		return me
	}

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
