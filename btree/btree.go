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

// Node
type Node struct {
	parent *Node
	left   *Node
	right  *Node
	value  Element
}

func (n *Node) Left() *Node {
	return n.left
}

func (n *Node) Right() *Node {
	return n.right
}

func (n *Node) setLeft(v *Node) {
	n.left = v
}

func (n *Node) setRight(v *Node) {
	n.right = v
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
	nextNode := func(prev, current *Node) *Node {
		if current == nil {
			log.Fatal("assertion failed! (current == nil)")
		}

		switch prev {
		case nil:
			return current.Left()
		case current.Left():
			if current.right != nil {
				return current.right
			} else {
				return current.parent
			}
		case current.right:
			if current.parent != nil {
				return current.parent
			} else {
				return nil
			}
		case current.parent:
			if current.Left() != nil {
				return current.Left()
			} else if current.right != nil {
				return current.right
			} else {
				return current.parent
			}
		default:
			log.Fatal("unreachable code!")
			return nil
		}
	}

	next := nextNode(prev, current)
	return current, next
}

func (btree *Btree) traversePrint() {
	logPrefix := log.Prefix()
	defer log.SetPrefix(logPrefix)
	log.SetPrefix("traversePrint: ")

	current := btree.root
	var prev *Node

	for current != nil {
		if prev == nil || prev == current.parent {
			log.Printf("value: %d\n", current.value)
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
			if current.right == nil {
				return current
			}
			current = current.right
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
			if current.right == nil {
				return nil
			}
			current = current.right
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
	newNode := &Node{
		parent: node,
		value:  v,
	}

	switch {
	case node == nil:
		btree.root = newNode
	case v.Lt(node.value):
		node.setLeft(newNode)
	case node.value.Lt(v):
		node.setRight(newNode)
	default:
		// ここに到達するとnewNodeは捨てられる
		return false
	}

	btree.len += 1
	return true
}

// Remove
func (node *Node) splice(btree *Btree) {
	if node.Left() != nil && node.right != nil {
		log.Fatal("unreachable")
	}

	var c *Node
	if node.Left() == nil {
		c = node.right
	} else {
		c = node.Left()
	}

	var p *Node

	if node == btree.root {
		btree.root = c
		p = nil
	} else {
		p = node.parent
		switch node {
		case p.Left():
			p.setLeft(c)
		case p.right:
			p.setRight(c)
		default:
			log.Fatal("unreachable")
		}
	}
	if c != nil {
		c.parent = p
	}
}

func (node *Node) remove(btree *Btree) {
	if node.Left() == nil || node.right == nil {
		node.splice(btree)
		return
	}

	alt := node.right
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

// recursivePrint
func (current *Node) recursivePrint() {
	log.Printf("%d\n", current.value)
	if current.Left() != nil {
		current.Left().recursivePrint()
	}
	if current.right != nil {
		current.right.recursivePrint()
	}
}

func (btree *Btree) recursivePrint() {
	logPrefix := log.Prefix()
	defer log.SetPrefix(logPrefix)
	log.SetPrefix("recursivePrint: ")

	btree.root.recursivePrint()
}

// traverseSetParent
func (current *Node) traverseSetParent(parent *Node) {
	if current == nil {
		return
	}

	current.parent = parent
	if current.Left() != nil {
		current.Left().traverseSetParent(current)
	}
	if current.right != nil {
		current.right.traverseSetParent(current)
	}
}

func (btree *Btree) traverseSetParent() {
	logPrefix := log.Prefix()
	defer log.SetPrefix(logPrefix)
	log.SetPrefix("traverseSetParent: ")

	btree.root.traverseSetParent(nil)
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
