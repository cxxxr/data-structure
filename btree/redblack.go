package btree

import (
	"errors"
	"fmt"
	"log"
)

type color int

const (
	red   color = 0
	black       = 1
)

type redBlackTree struct {
	root *redBlackNode
}

type redBlackNode struct {
	color  color
	value  Element
	parent *redBlackNode
	left   *redBlackNode
	right  *redBlackNode
}

func newRedBlackNode(color color, value Element, parent *redBlackNode) *redBlackNode {
	return &redBlackNode{
		color:  color,
		value:  value,
		parent: parent,
	}
}

func (n *redBlackNode) Left() Node {
	return n.left
}

func (n *redBlackNode) Right() Node {
	return n.right
}

func (n *redBlackNode) Value() Element {
	return n.value
}

func (n *redBlackNode) Color() color {
	if n == nil || n.color == black {
		return black
	}
	return red
}

func (n *redBlackNode) isLeftColor(color color) bool {
	return n.left.Color() == color
}

func (n *redBlackNode) isRightColor(color color) bool {
	return n.right.Color() == color
}

func (n *redBlackNode) pushBlack() {
	if !(n.color == black) {
		log.Fatal("assertion failed")
	}
	if !(n.isLeftColor(red)) {
		log.Fatal("assertion failed")
	}
	if !(n.isRightColor(red)) {
		log.Fatal("assertion failed")
	}

	n.color = red
	if n.left != nil {
		n.left.color = black
	}
	if n.right != nil {
		n.right.color = black
	}
}

func (n *redBlackNode) pullBlack() {
	if !(n.color == red) {
		log.Fatal("assertion failed")
	}
	if !(n.isLeftColor(black)) {
		log.Fatal("assertion failed")
	}
	if !(n.isRightColor(black)) {
		log.Fatal("assertion failed")
	}

	n.color = black
	if n.left != nil {
		n.left.color = red
	}
	if n.left != nil {
		n.right.color = red
	}
}

func (n *redBlackNode) swapColor(w *redBlackNode) {
	n.color, w.color = w.color, n.color
}

func (tree *redBlackTree) rotateLeft(n *redBlackNode) {
	r := n.right
	r.parent = n.parent
	n.parent = r
	n.right = r.left
	r.left = n
	if n == tree.root {
		tree.root = r
	}
}

func (tree *redBlackTree) rotateRight(n *redBlackNode) {
	l := n.left
	l.parent = n.parent
	n.parent = l
	n.left = l.right
	l.right = n
	if n == tree.root {
		tree.root = l
	}
}

func (tree *redBlackTree) flipLeft(n *redBlackNode) {
	n.swapColor(n.right)
	tree.rotateLeft(n)
}

func (tree *redBlackTree) flipRight(n *redBlackNode) {
	n.swapColor(n.left)
	tree.rotateRight(n)
}

func (tree *redBlackTree) addFixup(n *redBlackNode) {
	for n.color == red {
		if n == tree.root {
			n.color = black
			return
		}
		w := n.parent
		if w.isLeftColor(black) {
			tree.flipLeft(w)
			n = w
			w = n.parent
		}
		if w.color == black {
			return
		}
		g := w.parent
		if g.isRightColor(black) {
			log.Println("case-1")
			tree.flipRight(g)
			return
		} else {
			log.Println("case-2")
			g.pushBlack()
			n = g
		}
	}
}

func (tree *redBlackTree) findLastNode(v Element) *redBlackNode {
	current := tree.root
	var prev *redBlackNode

	for current != nil {
		prev = current
		if v.Eq(current.value) {
			return current
		}

		if v.Lt(current.value) {
			if current.left == nil {
				return current
			}
			current = current.left
		} else {
			if current.right == nil {
				return current
			}
			current = current.right
		}
	}

	return prev
}

func (tree *redBlackTree) Add(v Element) (*redBlackNode, error) {
	if tree == nil {
		return nil, errors.New("tree is nil")
	}

	n := tree.findLastNode(v)
	c := newRedBlackNode(red, v, n)

	switch {
	case n == nil:
		tree.root = c
	case v.Lt(n.value):
		n.left = c
	case n.value.Lt(v):
		n.right = c
	default:
		// Already has a value
		return nil, nil
	}

	tree.addFixup(c)

	return c, nil
}

func (node *redBlackNode) print(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
	switch node.color {
	case red:
		fmt.Print("r: ")
	case black:
		fmt.Print("b: ")
	}
	fmt.Printf("%v\n", node.value)
	if node.left != nil {
		fmt.Print("L; ")
		node.left.print(indent + 1)
	}
	if node.right != nil {
		fmt.Print("R; ")
		node.right.print(indent + 1)
	}
}

func (tree *redBlackTree) Print() {
	fmt.Print(" ; ")
	tree.root.print(0)
}
