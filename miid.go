package instanceid

import (
	"strconv"
	"strings"
	"time"
)

type StdMiid struct {
	sn string
	vn string
	va string
	t  int
}

func NewStdMiid(s string) *StdMiid {
	return parseMIID(s)
}

func (m StdMiid) Sn() string {
	return m.sn
}

func (m StdMiid) Vn() string {
	return m.vn
}

func (m StdMiid) Va() string {
	return m.va
}

func (m StdMiid) T() int {
	return m.t
}

func (m *StdMiid) SetT(t int) Miid {
	m.t = t
	return m
}

// String returns the textual representation of the Miid
func (m *StdMiid) String() string {
	sB := strings.Builder{}
	if m.sn != "" {
		sB.WriteString(m.sn)
		if m.vn != "" {
			sB.WriteString("/" + m.vn)
		}
		if m.va != "" {
			sB.WriteString("/" + m.va)
		}
		sB.WriteString("%" + strconv.Itoa(m.t) + "s")
	}
	return sB.String()
}

// Contains returns true if s is contained left aligned, else or if s is empty return false
func (m StdMiid) Contains(s string) bool {
	if s == "" {
		return false
	}
	return strings.Contains(m.String(), s)
}

// SetEpoch sets the epoch field based on a given StartTime. Chainable.
func (m *StdMiid) SetEpoch(startTime time.Time) Miid {
	epoch := time.Since(startTime)
	m.t = int(epoch.Seconds())
	return m
}

func parseMIID(id string) (miid *StdMiid) {
	miid = new(StdMiid)
	if !SanityCheck(id) {
		return miid
	}
	s := strings.SplitN(id, "/", -1)
	l := len(s)
	var r = new(StdMiid)
	if l >= 1 {
		r.sn = s[0]
	}
	if l == 2 {
		e := strings.Split(s[1], "%")
		r.vn = e[0]
		if len(e) > 1 {
			t, _ := strconv.Atoi(strings.Split(e[1], "s")[0])
			r.t = t
		}
	} else if l >= 2 {
		r.vn = s[1]
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
		r.va = va.String()

		t, _ := strconv.Atoi(strings.Split(e[1], "s")[0])
		r.t = t
	}

	return r
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
