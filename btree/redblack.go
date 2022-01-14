package btree

import (
	"errors"
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

func newRedBlackNode(color color, value Element) *redBlackNode {
	return &redBlackNode{
		color: color,
		value: value,
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

func (n *redBlackNode) pushBlack() {
	if !(n.color == black) {
		log.Fatal("assertion failed (n.color == black)")
	}
	if !(n.left.color == red) {
		log.Fatal("assertion failed (n.left.color == red)")
	}
	if !(n.right.color == red) {
		log.Fatal("assertion failed (n.right.color == red)")
	}
	n.color = red
	n.left.color = black
	n.right.color = black
}

func (n *redBlackNode) pullBlack() {
	if !(n.color == red) {
		log.Fatal("assertion failed (n.color == red")
	}
	if !(n.left.color == black) {
		log.Fatal("assertion failed (n.left.color == black")
	}
	if !(n.right.color == black) {
		log.Fatal("assertion failed (n.right.color == black")
	}
	n.color = black
	n.left.color = red
	n.right.color = red
}

func swapColor(n1, n2 *redBlackNode) {
	n1.color, n2.color = n2.color, n1.color
}

func (n *redBlackNode) rotateLeft() {
	r := n.right
	r.parent = n.parent
	n.parent = r
	n.right = r.left
	r.left = n
}

func (n *redBlackNode) rotateRight() {
	l := n.left
	l.parent = n.parent
	n.parent = l
	n.left = l.right
	l.right = n
}

func (n *redBlackNode) flipLeft() {
	swapColor(n, n.right)
	n.rotateLeft()
}

func (n *redBlackNode) flipRight() {
	swapColor(n, n.left)
	n.rotateRight()
}

func (tree *redBlackTree) addFixup(n *redBlackNode) {
	for n.color == red {
		if n == tree.root {
			n.color = black
			return
		}
		if n.parent.left.color == black {
			n.parent.flipLeft()
			n = n.parent
		}
		if n.parent.color == black {
			return
		}
		if n.parent.parent.right.color == black {
			n.parent.parent.flipRight()
			return
		} else {
			n.parent.parent.pushBlack()
			n = n.parent.parent
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
	c := newRedBlackNode(red, v)

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
