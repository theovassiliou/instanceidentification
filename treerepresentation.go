package instanceid

import (
	"strconv"

	"github.com/xlab/treeprint"
)

// PrintCiid prints a tree representation of the complete call-graph of a
// Ciid
func (ciid *Ciid) PrintCiid() string {
	tree := treeprint.New()
	tree = ciid.visitCiid(tree)
	return tree.String()
}

func (c *Ciid) visitCiid(t treeprint.Tree) treeprint.Tree {
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
	return m.String()
}

// PrintExtendedCiid prints a tree representation of the complete call-graph of a
// Ciid, and decodes external representations
func (ciid *Ciid) PrintExtendedCiid() string {
	tree := treeprint.New()
	tree = ciid.visitExtendedCiid(tree)
	return tree.String()
}

func (c *Ciid) visitExtendedCiid(t treeprint.Tree) treeprint.Tree {
	x := t.AddBranch(c.Miid.Sn + "/" + c.Miid.Vn)
	if c.Miid.Vn == "x" {
		x.SetMetaValue(c.Miid.decode(c.Miid.Va))
	} else if c.Miid.metadata() != "" {
		x.SetMetaValue(strconv.Itoa(c.Miid.T) + "s")
	}
	for _, s := range c.Ciids {
		s.visitExtendedCiid(x)
	}
	return t
}
