package instanceid

import (
	"strconv"

	"github.com/xlab/treeprint"
)

// TreePrint prints a tree representation of the complete call-graph of a
// Ciid
func (c StdCiid) TreePrint() string {
	tree := treeprint.New()
	tree = c.visitCiid(tree)
	return tree.String()
}

func (c StdCiid) visitCiid(t treeprint.Tree) treeprint.Tree {
	x := t.AddBranch(c.miid.Sn() + "/" + c.miid.Vn())
	if c.Miid().(StdMiid).metadata() != "" {
		x.SetMetaValue(strconv.Itoa(c.Miid().T()) + "s")
	}
	for _, s := range c.Ciids() {
		s.(StdCiid).visitCiid(x)
	}
	return t
}

func (m StdMiid) metadata() string {
	return m.String()
}
