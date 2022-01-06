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

type Node struct {
	edges []*Node
	value Element
}

func newNode(parent *Node, value Element) *Node {
	n := &Node{
		value: value,
		edges: make([]*Node, edgeNum),
	}
	n.edges[0] = parent
	return n
}

func (n *Node) Parent() *Node {
	return n.edges[edgeIndexParent]
}

func (n *Node) Left() *Node {
	return n.edges[edgeIndexLeft]
}

func (n *Node) Right() *Node {
	return n.edges[edgeIndexRight]
}

func (n *Node) setParent(v *Node) {
	n.edges[edgeIndexParent] = v
}

func (n *Node) setLeft(v *Node) {
	n.edges[edgeIndexLeft] = v
}

func (n *Node) setRight(v *Node) {
	n.edges[edgeIndexRight] = v
}

func (n *Node) children() []*Node {
	return n.edges[1:]
}

type Btree struct {
	root *Node
	len  int
}

func (btree *Btree) Len() int {
	return btree.len
}

func relativeEdgeIndex(prev, current *Node) edgeIndex {
	for edge := edgeIndexParent; edge < edgeNum; edge++ {
		if current.edges[edge] == prev {
			return edge
		}
	}
	panic("unreachable code")
}

func nextNode(prev, current *Node) (*Node, *Node) {
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

func (btree *Btree) Traverse(fn func(*Node)) {
	current := btree.root
	var prev *Node

	for current != nil {
		if prev == current.Parent() {
			fn(current)
		}
		prev, current = nextNode(prev, current)
	}
}

func (btree *Btree) findLastNode(v Element) *Node {
	current := btree.root
	var prev *Node

	for current != nil {
		prev = current
		if v.Eq(current.value) {
			return current
		}

		if v.Lt(current.value) {
			if current.Left() == nil {
				return current
			}
			current = current.Left()
		} else {
			if current.Right() == nil {
				return current
			}
			current = current.Right()
		}
	}

	return prev
}

func (btree *Btree) findNode(v Element) *Node {
	if btree == nil {
		return nil
	}

	current := btree.root
	for {
		if v.Eq(current.value) {
			return current
		}

		if v.Lt(current.value) {
			if current.Left() == nil {
				return nil
			}
			current = current.Left()
		} else {
			if current.Right() == nil {
				return nil
			}
			current = current.Right()
		}
	}
}

func (btree *Btree) Find(v Element) (*Node, error) {
	if btree == nil {
		return nil, errors.New("btree is nil")
	}
	node := btree.findNode(v)
	return node, nil
}

func (btree *Btree) Add(v Element) (*Node, error) {
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

func (node *Node) splice(btree *Btree) {
	if node.Left() != nil && node.Right() != nil {
		log.Fatal("unreachable")
	}

	var c *Node
	if node.Left() == nil {
		c = node.Right()
	} else {
		c = node.Left()
	}

	var p *Node

	if node == btree.root {
		btree.root = c
		p = nil
	} else {
		p = node.Parent()
		switch node {
		case p.Left():
			p.setLeft(c)
		case p.Right():
			p.setRight(c)
		default:
			log.Fatal("unreachable")
		}
	}
	if c != nil {
		c.setParent(p)
	}
}

func (node *Node) remove(btree *Btree) {
	if node.Left() == nil || node.Right() == nil {
		node.splice(btree)
		return
	}

	alt := node.Right()
	for alt.Left() != nil {
		alt = alt.Left()
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

func (node *Node) height() int {
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

	left := btree.root.Left()
	right := btree.root.Right()
	diff := left.height() - right.height()
	if 0 < diff {
		return diff < 2, nil
	} else {
		return -diff < 2, nil
	}
}
