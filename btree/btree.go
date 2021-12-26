package btree

import (
	"log"
	"fmt"
)

type Element interface{
	Eq(Element) bool
	Lt(Element) bool
	fmt.Stringer
}

type edgeIndex int

const (
	edgeIndexParent edgeIndex = iota
	edgeIndexLeft
	edgeIndexRight
	edgeNum
)

// Node
type Node struct {
	edges []*Node
	value  Element
}

func newNode(parent *Node, value Element) *Node {
	n := &Node{
		value: value,
		edges: make([]*Node, edgeNum),
	}
	n.edges[0] = parent
	return n
}

func (n *Node) Parent() *Node{
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

// Btree
type Btree struct {
	root *Node
	len  int
}

func (btree *Btree) Len() int {
	return btree.len
}

// traverse
func next(prev, current *Node) (*Node, *Node) {
	nextNode := func() *Node {
		if current == nil {
			log.Fatal("assertion failed! (current == nil)")
		}

		switch prev {
		case current.Left():
			if current.Right() != nil {
				return current.Right()
			} else {
				return current.Parent()
			}
		case current.Right():
			if current.Parent() != nil {
				return current.Parent()
			} else {
				return nil
			}
		case current.Parent():
			if current.Left() != nil {
				return current.Left()
			} else if current.Right() != nil {
				return current.Right()
			} else {
				return current.Parent()
			}
		default:
			log.Fatal("unreachable code!")
			return nil
		}
	}

	next := nextNode()
	return current, next
}

func (btree *Btree) Traverse(fn func(*Node)) {
	logPrefix := log.Prefix()
	defer log.SetPrefix(logPrefix)
	log.SetPrefix("traversePrint: ")

	current := btree.root
	prev := current.Parent() // rootのparentなのでnilになる

	for current != nil {
		if prev == current.Parent() {
			fn(current)
		}
		prev, current = next(prev, current)
	}
}

// Find
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

func (btree *Btree) Find(v Element) bool {
	node := btree.findNode(v)
	return node != nil
}

// Add
func (btree *Btree) Add(v Element) bool {
	if btree == nil {
		log.Fatal("assertion failed (btree.root == nil)")
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
		// ここに到達するとnewNodeは捨てられる
		return false
	}

	btree.len++
	return true
}

// Remove
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

func (btree *Btree) Remove(v Element) bool {
	if btree == nil {
		log.Fatal("assertion failed (btree.root == nil)")
	}

	node := btree.findNode(v)
	if node == nil {
		return false
	}

	node.remove(btree)
	btree.len--

	return true
}

// IntElement
type IntElement int

func (lhs IntElement) Eq(rhs Element) bool {
	v := rhs.(IntElement)
	return int(lhs) == int(v)
}

func (lhs IntElement) Lt(rhs Element) bool {
	v := rhs.(IntElement)
	return int(lhs) < int(v)
}

func (e IntElement) String() string {
	return fmt.Sprintf("%d", int(e))
}

// RuneElement
type RuneElement rune

func (lhs RuneElement) Eq(rhs Element) bool {
	v := rhs.(RuneElement)
	return rune(lhs) == rune(v)
}

func (lhs RuneElement) Lt(rhs Element) bool {
	v := rhs.(RuneElement)
	return rune(lhs) < rune(v)
}

func (e RuneElement) String() string {
	return string(rune(e))
}
