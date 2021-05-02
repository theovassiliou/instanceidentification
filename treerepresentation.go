package instanceid

import (
	"strconv"
	"strings"

	"github.com/xlab/treeprint"
)

// PrintCiid prints a tree representation of the complete call-graph of a
// Ciid
func PrintCiid(ciid Ciid) string {
	tree := treeprint.New()
	tree = ciid.visitCiid(tree)
	return tree.String()
}

func (c Ciid) visitCiid(t treeprint.Tree) treeprint.Tree {
	x := t.AddBranch(c.Miid.Sn + "/" + c.Miid.Vn)
	if c.Miid.metadata() != "" {
		x.SetMetaValue(strconv.Itoa(c.Miid.T) + "s")
	}
	for _, s := range c.Ciids {
		s.visitCiid(x)
	}
	return t
}

func (m *Miid) metadata() string {
	sB := strings.Builder{}

	if m.Vn != "" {
		sB.WriteString(m.Vn)
	}
	if m.Va != "" {
		sB.WriteString("/" + m.Va)
	}
	if m.T != 0 {
		sB.WriteString("%" + strconv.Itoa(m.T) + "s")
	}
	return sB.String()
}
