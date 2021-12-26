package btree

import (
	"fmt"
	"strings"
	"testing"
)

func newBtree() *Btree {
	btree := &Btree{
		root: &Node{
			left: &Node{
				value: 1,
			},
			value: 2,
			right: &Node{
				value: 3,
			},
		},
		len: 3,
	}
	btree.traverseSetParent()
	return btree
}

func TestFind(t *testing.T) {
	btree := newBtree()
	if !btree.Find(1) {
		t.Fatal("unexpected: btree.Find(1) is false")
	}
	if !btree.Find(2) {
		t.Fatal("unexpected: btree.Find(2) is false")
	}
	if !btree.Find(3) {
		t.Fatal("unexpected: btree.Find(3) is false")
	}
	if btree.Find(4) {
		t.Fatal("unexpected: btree.Find(4) is true")
	}
}

func TestPrint(t *testing.T) {
	btree := newBtree()
	btree.recursivePrint()
	btree.traversePrint()
}

func TestAdd(t *testing.T) {
	var btree Btree
	values := []int{1, 7, 4, 0, 9, 2, 3, 5, 8, 6}
	for i, v := range values {
		if !btree.Add(v) {
			t.Fatalf("!btree.Add(%v)", v)
		}
		if !btree.Find(v) {
			t.Fatalf("!btree.Find(%v)", v)
		}
		if btree.Len() != i+1 {
			t.Fatalf("btree.Len() != %d\n", i+1)
		}
	}

	for _, v := range values {
		if btree.Add(v) {
			t.Fatalf("btree.Add(%v)", v)
		}
	}
}

func genDotText(graph [][]int) string {
	var b strings.Builder
	b.WriteString("digraph btree {\n")
	for _, vec := range graph {
		from := vec[0]
		to := vec[1]
		b.WriteString(fmt.Sprintf("%d -> %d;\n", from, to))
	}
	b.WriteString("}\n")
	return b.String()
}

func genTestingBtree(values []int) *Btree {
	var btree Btree
	for _, v := range values {
		btree.Add(v)
	}
	return &btree
}

func toDot(btree *Btree) string {
	var builder strings.Builder
	btree.PrintDot(&builder)
	return builder.String()
}

func testBtreeShape(t *testing.T, btree *Btree, expected string) {
	actual := toDot(btree)
	if expected != actual {
		t.Fatalf("expected = %v\nactual = %v", expected, actual)
	}
	if btree.Len() != strings.Count(expected, ";")+1 {
		t.Fatal("unexpected Len() values")
	}
}

func TestRemove(t *testing.T) {
	var expected string

	values := []int{7, 3, 11, 1, 5, 9, 13, 4, 6, 8, 12, 14}
	btree := genTestingBtree(values)

	btree.Remove(6)
	expected = genDotText([][]int{
		{7, 3},
		{3, 1},
		{3, 5},
		{5, 4},
		{7, 11},
		{11, 9},
		{9, 8},
		{11, 13},
		{13, 12},
		{13, 14},
	})
	testBtreeShape(t, btree, expected)

	btree.Remove(9)
	expected = genDotText([][]int{
		{7, 3},
		{3, 1},
		{3, 5},
		{5, 4},
		{7, 11},
		{11, 8},
		{11, 13},
		{13, 12},
		{13, 14},
	})
	testBtreeShape(t, btree, expected)

	btree.Remove(11)
	expected = genDotText([][]int{
		{7, 3},
		{3, 1},
		{3, 5},
		{5, 4},
		{7, 12},
		{12, 8},
		{12, 13},
		{13, 14},
	})
	testBtreeShape(t, btree, expected)

	btree.Remove(7)
	expected = genDotText([][]int{
		{8, 3},
		{3, 1},
		{3, 5},
		{5, 4},
		{8, 12},
		{12, 13},
		{13, 14},
	})
	testBtreeShape(t, btree, expected)

	btree = genTestingBtree([]int{1, 2, 3})
	btree.Remove(1)
	expected = genDotText([][]int{
		{2, 3},
	})
	testBtreeShape(t, btree, expected)
}
