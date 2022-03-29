package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var printFiles bool

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles = len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func walkFunc(out io.Writer, path string, prefix string) error {
	contents, _ := os.ReadDir(path)

	for i, c := range contents {
		if strings.HasPrefix(c.Name(), ".") {
			continue
		}
		if c.IsDir() {
			if i == len(contents)-1 {
				fmt.Fprintln(out, prefix+"└───"+c.Name())
				walkFunc(out, path+"/"+c.Name(), prefix+"\t")
			} else {
				fmt.Fprintln(out, prefix+"├───"+c.Name())
				walkFunc(out, path+"/"+c.Name(), prefix+"│\t")
			}
		} else {
			fi, err := os.Stat(path + "/" + c.Name())
			if err != nil {
				log.Fatal(err)
			}
			if i == len(contents)-1 {
				fmt.Fprintln(out, fmt.Sprintf("%s└───%s (%vb)", prefix, c.Name(), fi.Size()))
			} else {
				fmt.Fprintln(out, fmt.Sprintf("%s├───%s (%vb)", prefix, c.Name(), fi.Size()))
			}
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	fmt.Fprintln(out, printFiles)

	if err := walkFunc(out, path, ""); err != nil {
		fmt.Printf("Какая-то ошибка: %v\n", err)
	}

	return nil
}
