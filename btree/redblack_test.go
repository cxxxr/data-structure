package btree

import (
	"fmt"
	"testing"
)

func TestRedBlack(t *testing.T) {
	var tree redBlackTree
	for i := 1; i <= 3; i++ {
		fmt.Printf("### %v ###\n", i)
		tree.Add(Int(i))
		tree.Print()
	}
	// GenDotAndOpenImage("_redblack", tree.root)
}
