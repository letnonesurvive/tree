package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

const beginSymbol string = "├"
const lastDirSymbol string = "└"

type visitorTree struct {
	out     io.Writer
	depth   int
	pipes   int
	builder strings.Builder
}

func printTab(visitor *visitorTree) {
	for j := 0; j < visitor.depth; j++ {
		if j < visitor.pipes {
			visitor.builder.WriteString("│\t")
		} else {
			visitor.builder.WriteString("\t")
		}
	}
}

func (visitor *visitorTree) visitEnter(tree *treeNode) {

	if len(tree.children) != 0 {
		visitor.depth++
	}

	for i := 0; i < len(tree.children); i++ {
		printTab(visitor)
		file := tree.children[i].value
		if i < len(tree.children)-1 {
			visitor.builder.WriteString(beginSymbol)
			if len(tree.children[i].children) != 0 {
				visitor.pipes++
			}
		} else {
			visitor.builder.WriteString(lastDirSymbol)
		}
		fmt.Fprintf(&visitor.builder, "───%s", file.name)
		if file.fileType == RegularFile {
			fmt.Fprintf(&visitor.builder, " (%s)", file.size)
		}
		visitor.builder.WriteString("\n")
		if len(tree.children[i].children) != 0 {
			visitor.visitEnter(tree.children[i])
			visitor.visitLeave(tree.children[i])
		}
	}
}

func (visitor *visitorTree) visitLeave(_ *treeNode) {
	if visitor.depth != 0 {
		visitor.depth--
	}
	if visitor.pipes != 0 && visitor.pipes != visitor.depth {
		visitor.pipes--
	}
}

type treeNode struct {
	value    file
	children []*treeNode
}

type fileType uint

const (
	RegularFile = iota
	Dirrectory
)

type file struct {
	name     string
	size     string
	fileType fileType
}

func (tree *treeNode) accept(visitor *visitorTree) {
	visitor.visitEnter(tree)
	visitor.visitLeave(tree)
}

func buildTree(path string, printFiles bool) []*treeNode {
	currentEntries, _ := os.ReadDir(path)

	var files []file

	if !printFiles {
		for i := 0; i < len(currentEntries); i++ {
			if !currentEntries[i].IsDir() {
				continue
			}
			var file file
			file.fileType = Dirrectory
			file.name = currentEntries[i].Name()
			files = append(files, file)
		}
	} else {
		for i := 0; i < len(currentEntries); i++ {
			size := ""
			var file file
			file.fileType = Dirrectory
			file.name = currentEntries[i].Name()
			var entry fs.DirEntry = currentEntries[i]
			info, _ := entry.Info()
			if !entry.IsDir() {
				file.fileType = RegularFile
				size = "empty"
				if info.Size() != 0 {
					size = strconv.Itoa(int(info.Size())) + "b"
				}
				file.size = size
			}

			files = append(files, file)
		}
	}

	res := make([]*treeNode, len(files))

	for i := 0; i < len(files); i++ {
		treeNode := new(treeNode)
		treeNode.value = files[i]
		treeNode.children = buildTree(path+"/"+treeNode.value.name, printFiles)
		res[i] = treeNode
	}
	return res
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	root := new(treeNode)
	var file file
	root.value = file
	root.children = buildTree(path, printFiles)

	var visitor visitorTree
	visitor.out = out
	visitor.depth = -1
	visitor.pipes = 0
	visitor.builder = strings.Builder{}
	root.accept(&visitor)

	fmt.Fprintf(out, visitor.builder.String())

	return nil
}

func main() {
	//out, _ := os.Create("output.txt")
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
