package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func printDir(out io.Writer, n int, dirName string, isLastDir bool, isParentDirLast bool) {
	beginSymbol := "├"
	lastDirSymbol := "└"
	for i := 0; i < n-1; i++ {
		if isParentDirLast || i == n-2 {
			fmt.Fprintf(out, "|")
		}
		fmt.Fprintf(out, "\t")
	}
	if isLastDir {
		fmt.Fprintf(out, lastDirSymbol)
	} else {
		fmt.Fprintf(out, beginSymbol)
	}
	fmt.Fprintf(out, "───"+dirName+"\n")

}

// func printDirs(out io.Writer, dirs []string) {
// 	for i := 0; i < len(dirs); i++ {
// 		nesting := strings.SplitN(dirs[i], "/", -1)
// 		n := len(nesting)
// 		isLastDir := (i == len(dirs)-1) || (strings.Count(dirs[i+1], "/") < n-1)
// 		printDir(out, n, nesting[n-1], isLastDir)
// 	}
// }

func getFiles(path string, printFiles bool) []string {

	// entries, _ := os.ReadDir(path)
	// files := make(map[string]bool, len(entries))

	// sort.Slice(entries, func(i, j int) bool {
	// 	return entries[i].Name() < entries[j].Name()
	// })

	// for i := 0; i < len(entries); i++ {
	// 	if printFiles {
	// 		files[entries[i].Name()] = false
	// 	} else if !printFiles && entries[i].IsDir() {
	// 		files[entries[i].Name()] = false
	// 	}
	// }

	//fmt.Println(files)
	var dirs []string

	filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		entries, _ := os.ReadDir(path)
		fmt.Println(entries[len(entries)-1])
		// if printFiles {
		// 	if strings.Contains(p, path+"/") {
		// 		p = strings.Replace(p, path+"/", "", 1)
		// 		if d.Type().IsRegular() {
		// 			var size string = ""
		// 			info, _ := d.Info()
		// 			size = strconv.FormatInt(info.Size(), 10)
		// 			if size == "0" {
		// 				size = "empty"
		// 			}
		// 			p += " (" + size + ")"
		// 		}
		// 		dirs = append(dirs, p)
		// 	}
		// } else {
		if strings.Contains(p, path+"/") && d.IsDir() {
			p = strings.Replace(p, path+"/", "", 1)
			dirs = append(dirs, p)
		}
		//}
		return nil
	})
	return dirs
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	dirs := getFiles(path, printFiles)
	//printDirs(out, dirs)
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	//printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, false)
	if err != nil {
		panic(err.Error())
	}
}
