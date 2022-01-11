package btree

import (
	"errors"
	"log"
)

type edgeIndex int

const (
	edgeIndexParent edgeIndex = iota
	edgeIndexLeft
	edgeIndexRight
	edgeNum
)

func nextEdge(edge edgeIndex) edgeIndex {
	return (edge + 1) % edgeNum
}

type BinaryNode struct {
	edges []*BinaryNode
	value Element
}

func newNode(parent *BinaryNode, value Element) *BinaryNode {
	n := &BinaryNode{
		value: value,
		edges: make([]*BinaryNode, edgeNum),
	}
	n.edges[0] = parent
	return n
}

func (n *BinaryNode) parent() *BinaryNode {
	return n.edges[edgeIndexParent]
}

func (n *BinaryNode) left() *BinaryNode {
	return n.edges[edgeIndexLeft]
}

func (n *BinaryNode) right() *BinaryNode {
	return n.edges[edgeIndexRight]
}

func (n *BinaryNode) setParent(v *BinaryNode) {
	n.edges[edgeIndexParent] = v
}

func (n *BinaryNode) setLeft(v *BinaryNode) {
	n.edges[edgeIndexLeft] = v
}

func (n *BinaryNode) setRight(v *BinaryNode) {
	n.edges[edgeIndexRight] = v
}

func (n *BinaryNode) children() []*BinaryNode {
	return n.edges[1:]
}

type Btree struct {
	root *BinaryNode
	len  int
}

func (btree *Btree) Len() int {
	return btree.len
}

func relativeEdgeIndex(prev, current *BinaryNode) edgeIndex {
	for edge := edgeIndexParent; edge < edgeNum; edge++ {
		if current.edges[edge] == prev {
			return edge
		}
	}
	panic("unreachable code")
}

func nextNode(prev, current *BinaryNode) (*BinaryNode, *BinaryNode) {
	// 以下の順でnilではないedgeを探す
	// parent : left -> right -> parent
	// left   : right -> parent
	// right  : parent

	edge := relativeEdgeIndex(prev, current)
	for {
		edge = nextEdge(edge)
		if n := current.edges[edge]; n != nil {
			return current, n
		}
		if edge == edgeIndexParent {
			return current, nil
		}
	}
}

func (btree *Btree) Traverse(fn func(*BinaryNode)) {
	current := btree.root
	var prev *BinaryNode

	for current != nil {
		if prev == current.parent() {
			fn(current)
		}
		prev, current = nextNode(prev, current)
	}
}

func (btree *Btree) findLastNode(v Element) *BinaryNode {
	current := btree.root
	var prev *BinaryNode

	for current != nil {
		prev = current
		if v.Eq(current.value) {
			return current
		}

		if v.Lt(current.value) {
			if current.left() == nil {
				return current
			}
			current = current.left()
		} else {
			if current.right() == nil {
				return current
			}
			current = current.right()
		}
	}

	return prev
}

func (btree *Btree) findNode(v Element) *BinaryNode {
	if btree == nil {
		return nil
	}

	current := btree.root
	for {
		if v.Eq(current.value) {
			return current
		}

		if v.Lt(current.value) {
			if current.left() == nil {
				return nil
			}
			current = current.left()
		} else {
			if current.right() == nil {
				return nil
			}
			current = current.right()
		}
	}
}

func (btree *Btree) Find(v Element) (*BinaryNode, error) {
	if btree == nil {
		return nil, errors.New("btree is nil")
	}
	node := btree.findNode(v)
	return node, nil
}

func (btree *Btree) Add(v Element) (*BinaryNode, error) {
	if btree == nil {
		return nil, errors.New("btree is nil")
	}

	node := btree.findLastNode(v)
	child := newNode(node, v)

	switch {
	case node == nil:
		btree.root = child
	case v.Lt(node.value):
		node.setLeft(child)
	case node.value.Lt(v):
		node.setRight(child)
	default:
		// 既に要素があるとここに到達する、newNodeは捨てられる
		return nil, nil
	}

	btree.len++
	return child, nil
}

func (node *BinaryNode) splice(btree *Btree) {
	if node.left() != nil && node.right() != nil {
		log.Fatal("unreachable")
	}

	var c *BinaryNode
	if node.left() == nil {
		c = node.right()
	} else {
		c = node.left()
	}

	var p *BinaryNode

	if node == btree.root {
		btree.root = c
		p = nil
	} else {
		p = node.parent()
		switch node {
		case p.left():
			p.setLeft(c)
		case p.right():
			p.setRight(c)
		default:
			log.Fatal("unreachable")
		}
	}
	if c != nil {
		c.setParent(p)
	}
}

func (node *BinaryNode) remove(btree *Btree) {
	if node.left() == nil || node.right() == nil {
		node.splice(btree)
		return
	}

	alt := node.right()
	for alt.left() != nil {
		alt = alt.left()
	}

	node.value = alt.value
	alt.splice(btree)
}

func (btree *Btree) Remove(v Element) (bool, error) {
	if btree == nil {
		return false, errors.New("btree is nil")
	}

	node := btree.findNode(v)
	if node == nil {
		return false, nil
	}

	node.remove(btree)
	btree.len--

	return true, nil
}

func (node *BinaryNode) height() int {
	if node == nil {
		return 0
	}

	max := 0
	for _, n := range node.children() {
		v := n.height()
		if max < v {
			max = v
		}
	}
	return 1 + max
}

func (btree *Btree) Height() (int, error) {
	if btree == nil {
		return 0, errors.New("btree is nil")
	}

	return btree.root.height(), nil
}

func (btree *Btree) Balanced() (bool, error) {
	if btree == nil {
		return false, errors.New("btree is nil")
	}

	left := btree.root.left()
	right := btree.root.right()
	diff := left.height() - right.height()
	if 0 < diff {
		return diff < 2, nil
	} else {
		return -diff < 2, nil
	}
}
