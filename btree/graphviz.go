package btree

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

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
