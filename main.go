package main

import (
	"github.com/cxxxr/btree/btree"
)

func main() {
	values := []btree.Int{7, 3, 11, 1, 5, 9, 13, 4, 6, 8, 12, 14}

	var btree btree.Btree
	for _, v := range values {
		btree.Add(v)
	}
}
