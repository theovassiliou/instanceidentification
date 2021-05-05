package instanceid

import (
	"strconv"
	"strings"
	"time"
)

// DecodingFunction of function type
type DecodingFunction func(string) string

// Miid represents the instance only by it's name, version, additional information
// and epoch time
type Miid struct {
	Sn string
	Vn string
	Va string
	T  int

	decode DecodingFunction
}

// NewMiid creates a new Miid from a string in the of
// Sn1/Vn1/Va1%t1s
// in case a Ciid is being provided the Miid part is only
// returned
// If there are syntax errors an empty Miid will be returned
func NewMiid(id string) (miid Miid) {
	return parseMIID(id)
}

func (c *Miid) WithDecoding(d DecodingFunction) *Miid {
	c.decode = d
	return c
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

// SetEpoch sets the epoch field based on a given StartTime. Chainable.
func (miid *Miid) SetEpoch(startTime time.Time) *Miid {
	epoch := time.Since(startTime)
	miid.T = int(epoch.Seconds())
	return miid
}

// String returns the textual representation of the Miid
func (m *Miid) String() string {
	sB := strings.Builder{}
	if m.Sn != "" {
		sB.WriteString(m.Sn)
		if m.Vn != "" {
			sB.WriteString("/" + m.Vn)
		}
		if m.Va != "" {
			if m.decode != nil {
				sB.WriteString("/" + m.decode(m.Va))
			} else {
				sB.WriteString("/" + m.Va)
			}
		}
		sB.WriteString("%" + strconv.Itoa(m.T) + "s")
	}
	return sB.String()
}

// Contains returns true if s is contained left aligned, else or if s is empty return false
func (m *Miid) Contains(s string) bool {
	if s == "" {
		return false
	}
	return strings.Contains(m.String(), s)
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

	if l >= 3 {
		e := strings.Split(s[len(s)-1], "%")

		rest := s[2:]
		va := strings.Builder{}
		for i := 0; i < len(rest)-1; i++ {
			va.WriteString(rest[i])
			va.WriteString("/")
		}
		va.WriteString(e[0])
		r.Va = va.String()

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
