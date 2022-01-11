package btree

import (
	"fmt"
	"strings"
	"testing"
)

func TestAddFind(t *testing.T) {
	var btree Btree
	values := []Int{1, 7, 4, 0, 9, 2, 3, 5, 8, 6}
	for i, v := range values {
		if n, err := btree.Add(v); !(n != nil && err == nil) {
			t.Fatalf("expected btree.Add(%v) to be (!nil, nil)", v)
		}
		if n, err := btree.Find(v); !(n != nil && err == nil) {
			t.Fatalf("expected btree.Find(%v) to be (!nil, nil)", v)
		}
		if btree.Len() != i+1 {
			t.Fatalf("btree.Len() != %d\n", i+1)
		}
	}

	for _, v := range values {
		if n, err := btree.Add(v); !(n == nil && err == nil) {
			t.Fatalf("expected btree.Add(%v) to be (nil, nil)", v)
		}
	}

	if n, err := btree.Find(Int(100)); !(n == nil && err == nil) {
		t.Fatalf("expected btree.Find(100) to be (nil, nil)")
	}
}

func TestFindNil(t *testing.T) {
	var btree *Btree
	n, err := btree.Find(Int(1))
	if !(n == nil && err != nil) {
		t.Fatal("at btree == nil, Find was expected to be (nil, !nil)")
	}
}

func TestAddNil(t *testing.T) {
	var btree *Btree
	n, err := btree.Add(Int(1))
	if !(n == nil && err != nil) {
		t.Fatal("at btree == nil, Add was expected to be (nil, !nil)")
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

func genTestingBtree(values []Int) *Btree {
	var btree Btree
	for _, v := range values {
		btree.Add(v)
	}
	return &btree
}

func toDot(btree *Btree) string {
	var builder strings.Builder
	GenDot(&builder, btree.root)
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
	values := []Int{7, 3, 11, 1, 5, 9, 13, 4, 6, 8, 12, 14}
	btree := genTestingBtree(values)

	if b, err := btree.Remove(Int(6)); !(b && err == nil) {
		t.Fatal("expected to be (true, nil)")
	}
	expected := genDotText([][]int{
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

	if b, err := btree.Remove(Int(9)); !(b && err == nil) {
		t.Fatal("expected to be (true, nil)")
	}
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

	if b, err := btree.Remove(Int(11)); !(b && err == nil) {
		t.Fatal("expected to be (true, nil)")
	}
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

	if b, err := btree.Remove(Int(7)); !(b && err == nil) {
		t.Fatal("expected to be (true, nil)")
	}
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

	btree = genTestingBtree([]Int{1, 2, 3})
	btree.Remove(Int(1))
	expected = genDotText([][]int{
		{2, 3},
	})
	testBtreeShape(t, btree, expected)
}

func TestRemoveNil(t *testing.T) {
	var btree *Btree
	b, err := btree.Remove(Int(1))
	if !(!b && err != nil) {
		t.Fatal("at btree == nil, Remove was expected to be (nil, !nil)")
	}
}

func TestTraverse(t *testing.T) {
	values := []Int{7, 3, 11, 1, 5, 9, 13, 4, 6, 8, 12, 14}
	btree := genTestingBtree(values)

	actual := make([]int, 0)
	btree.Traverse(func(n *BinaryNode) {
		actual = append(actual, int(n.value.(Int)))
	})

	if len(actual) != len(values) {
		t.Fatal("len(vec) != len(values)")
	}

	expected := []int{7, 3, 1, 5, 4, 6, 11, 9, 8, 13, 12, 14}
	for i := range actual {
		if actual[i] != int(expected[i]) {
			t.Fatalf("actual[%d] != expected[%d]", i, i)
		}
	}
}

func TestHeight(t *testing.T) {
	var btree *Btree
	if n, err := btree.Height(); !(n == 0 && err != nil) {
		t.Fatal("at btree == nil, Height was expected to be (0, !nil)")
	}

	btree = &Btree{}
	if n, err := btree.Height(); !(n == 0 && err == nil) {
		t.Fatal("expected btree.Height() to be 0")
	}

	values := []Int{7, 3, 11, 1, 5, 9, 13, 4, 6, 8, 12, 14}
	btree = genTestingBtree(values)
	if height, err := btree.Height(); !(height == 4 && err == nil) {
		t.Fatal("expected btree.Height() to be 4")
	}
}

func testBalanced(t *testing.T, btree *Btree, expected bool) {
	if balanced, err := btree.Balanced(); !(balanced == expected && err == nil) {
		t.Fatalf("expected btree.Balanced() == %v", expected)
	}
}

func TestBalanced(t *testing.T) {
	var btree *Btree

	if balanced, err := btree.Balanced(); !(balanced == false && err != nil) {
		t.Fatal("at btree == nil, Balanced was expected to be (false, !nil)")
	}

	btree = &Btree{}
	btree.Add(Int(1))
	testBalanced(t, btree, true)

	btree.Add(Int(2))
	testBalanced(t, btree, true)

	btree.Add(Int(3))
	testBalanced(t, btree, false)
}
