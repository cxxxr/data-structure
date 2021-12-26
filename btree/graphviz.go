package btree

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func (node *Node) dot(output io.StringWriter) {
	if node.Left() != nil {
		output.WriteString(fmt.Sprintf("%v -> %v;\n", node.value, node.Left().value))
		node.Left().dot(output)
	}
	if node.Right() != nil {
		output.WriteString(fmt.Sprintf("%v -> %v;\n", node.value, node.Right().value))
		node.Right().dot(output)
	}
}

func (btree *Btree) GenDot(output io.StringWriter) {
	output.WriteString("digraph btree {\n")
	btree.root.dot(output)
	output.WriteString("}\n")
}

func (btree *Btree) GenDotAndOpenImage(baseName string) {
	dotName := fmt.Sprintf("%s.dot", baseName)
	pngName := fmt.Sprintf("%s.png", baseName)

	output, err := os.Create(dotName)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	btree.GenDot(output)

	if err := exec.Command("dot", "-T", "png", dotName, "-o", pngName).Run(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("open", pngName).Run(); err != nil {
		log.Fatal(err)
	}
}
