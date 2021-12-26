package btree

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
)

// Node
type Node struct {
	left   *Node
	right  *Node
	parent *Node
	value  int
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
			return current.left
		case current.left:
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
			if current.left != nil {
				return current.left
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
func (btree *Btree) findLastNode(v int) *Node {
	current := btree.root
	var prev *Node

	for current != nil {
		prev = current
		if current.value == v {
			return current
		}

		if current.value < v {
			if current.right == nil {
				return current
			}
			current = current.right
		} else {
			if current.left == nil {
				return current
			}
			current = current.left
		}
	}

	return prev
}

func (btree *Btree) findNode(v int) *Node {
	if btree == nil {
		return nil
	}

	current := btree.root
	for {
		if current.value == v {
			return current
		}

		if current.value < v {
			if current.right == nil {
				return nil
			}
			current = current.right
		} else {
			if current.left == nil {
				return nil
			}
			current = current.left
		}
	}
}

func (btree *Btree) Find(v int) bool {
	node := btree.findNode(v)
	return node != nil
}

// Add
func (btree *Btree) Add(v int) bool {
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
	case node.value < v:
		node.right = newNode
	case node.value > v:
		node.left = newNode
	default:
		// ここに到達するとnewNodeは捨てられる
		return false
	}

	btree.len += 1
	return true
}

// Remove
func (node *Node) splice(btree *Btree) {
	if node.left != nil && node.right != nil {
		log.Fatal("unreachable")
	}

	var c *Node
	if node.left == nil {
		c = node.right
	} else {
		c = node.left
	}

	var p *Node

	if node == btree.root {
		btree.root = c
		p = nil
	} else {
		p = node.parent
		switch node {
		case p.left:
			p.left = c
		case p.right:
			p.right = c
		default:
			log.Fatal("unreachable")
		}
	}
	if c != nil {
		c.parent = p
	}
}

func (node *Node) remove(btree *Btree) {
	if node.left == nil || node.right == nil {
		node.splice(btree)
		return
	}

	alt := node.right
	for alt.left != nil {
		alt = alt.left
	}

	node.value = alt.value
	alt.splice(btree)
}

func (btree *Btree) Remove(v int) bool {
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
	if current.left != nil {
		current.left.recursivePrint()
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
	if current.left != nil {
		current.left.traverseSetParent(current)
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

// graphviz
func (node *Node) dot(output io.StringWriter) {
	if node.left != nil {
		output.WriteString(fmt.Sprintf("%d -> %d;\n", node.value, node.left.value))
		node.left.dot(output)
	}
	if node.right != nil {
		output.WriteString(fmt.Sprintf("%d -> %d;\n", node.value, node.right.value))
		node.right.dot(output)
	}
}

func (btree *Btree) PrintDot(output io.StringWriter) {
	output.WriteString("digraph btree {\n")
	btree.root.dot(output)
	output.WriteString("}\n")
}

func (btree *Btree) PrintDotAndOpenImage(baseName string) {
	dotName := fmt.Sprintf("%s.dot", baseName)
	pngName := fmt.Sprintf("%s.png", baseName)

	output, err := os.Create(dotName)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	btree.PrintDot(output)

	if err := exec.Command("dot", "-T", "png", dotName, "-o", pngName).Run(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("open", pngName).Run(); err != nil {
		log.Fatal(err)
	}
}

// main
func NewRandomValues() []int {
	values := make([]int, 10)
	for i := 0; i < 10; i++ {
		values[i] = i
	}

	rand.Shuffle(10, func(i, j int) {
		tmp := values[i]
		values[i] = values[j]
		values[j] = tmp
	})
	return values
}

/*
*/