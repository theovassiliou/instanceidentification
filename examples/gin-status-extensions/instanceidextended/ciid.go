package instanceidextended

import (
	"strings"

	instanceid "github.com/theovassiliou/instanceidentification"
)

type ExtendedCiid struct {
	instanceid.StdCiid
}

func NewExtCiid(a, v, b, c string) instanceid.Ciid {
	branchB := strings.Builder{}
	if b != "" || c != "" {
		branchB.WriteString("/")
	}
	if b != "" {
		branchB.WriteString(b)
	}

	if c != "" {
		branchB.WriteString("-")
		branchB.WriteString(c)
	}

	return instanceid.NewStdCiid(a + "/" + v + branchB.String() + "%-1s")
}
