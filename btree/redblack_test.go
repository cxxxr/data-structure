package btree

import "testing"

func TestRedBlack(t *testing.T) {
	var tree redBlackTree
	tree.Add(Int(1))
	tree.Add(Int(2))
	// GenDotAndOpenImage("_redblack.dot", tree.root)
}
