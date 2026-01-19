package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

const fullConnector = "├───"
const lastConnector = "└───"
const emptySign = " (empty)"
const tab = "\t"
const halfTab = "│\t"

func dirTree(out io.Writer, path string, printFiles bool) error {
	return walk(out, path, printFiles, []bool{})
}

func walk(out io.Writer, path string, printFiles bool, lastLevels []bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	filtered := make([]os.DirEntry, 0)
	for _, e := range entries {
		if e.Name() == ".DS_Store" {
			continue
		}
		if e.IsDir() || printFiles {
			filtered = append(filtered, e)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name() < filtered[j].Name()
	})

	for i, entry := range filtered {
		isLast := i == len(filtered)-1

		prefix := ""
		for _, last := range lastLevels {
			if last {
				prefix += tab
			} else {
				prefix += halfTab
			}
		}

		connector := fullConnector
		if isLast {
			connector = lastConnector
		}

		line := prefix + connector + entry.Name()

		if !entry.IsDir() {
			info, _ := entry.Info()
			if info.Size() == 0 {
				line += emptySign
			} else {
				line += fmt.Sprintf(" (%db)", info.Size())
			}
		}

		fmt.Fprintln(out, line)

		if entry.IsDir() {
			err := walk(out, filepath.Join(path, entry.Name()), printFiles, append(lastLevels, isLast))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
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
