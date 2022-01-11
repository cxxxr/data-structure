package btree

type Node interface {
	Left() Node
	Right() Node
	Value() Element
}
