package btree

import "fmt"

type Element interface {
	Eq(Element) bool
	Lt(Element) bool
	fmt.Stringer
}

type Int int

func (lhs Int) Eq(rhs Element) bool {
	v := rhs.(Int)
	return int(lhs) == int(v)
}

func (lhs Int) Lt(rhs Element) bool {
	v := rhs.(Int)
	return int(lhs) < int(v)
}

func (e Int) String() string {
	return fmt.Sprintf("%d", int(e))
}

type Rune rune

func (lhs Rune) Eq(rhs Element) bool {
	v := rhs.(Rune)
	return rune(lhs) == rune(v)
}

func (lhs Rune) Lt(rhs Element) bool {
	v := rhs.(Rune)
	return rune(lhs) < rune(v)
}

func (e Rune) String() string {
	return string(rune(e))
}
