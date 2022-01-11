package btree

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
)

func dot(output io.StringWriter, node Node) {
	if !reflect.ValueOf(node.Left()).IsNil() {
		output.WriteString(fmt.Sprintf("%v -> %v;\n", node.Value(), node.Left().Value()))
		dot(output, node.Left())
	}
	if !reflect.ValueOf(node.Right()).IsNil() {
		output.WriteString(fmt.Sprintf("%v -> %v;\n", node.Value(), node.Right().Value()))
		dot(output, node.Right())
	}
}

func GenDot(output io.StringWriter, root Node) {
	output.WriteString("digraph btree {\n")
	dot(output, root)
	output.WriteString("}\n")
}

func GenDotAndOpenImage(baseName string, root Node) {
	dotName := fmt.Sprintf("%s.dot", baseName)
	pngName := fmt.Sprintf("%s.png", baseName)

	output, err := os.Create(dotName)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	GenDot(output, root)

	if err := exec.Command("dot", "-T", "png", dotName, "-o", pngName).Run(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("open", pngName).Run(); err != nil {
		log.Fatal(err)
	}
}
