package main

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"log"
)

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

func dirTree(out io.Writer, filePath string, printFiles bool) error  {
	var dirGraph = ""
	err := levelDir(out, filePath, printFiles, dirGraph)

	return err
}

//Рекурсивная ф-ия обрабатывающая отдельный уровень дерева файлов
func levelDir(out io.Writer, filePath string, printFiles bool, dirGraph string) error {

	var tempGraph string
	var tempSize string
	var printGraph string

	files, err := ioutil.ReadDir(filePath)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
		if !file.IsDir() && !printFiles {
			continue
		}

		if file.Name() == ".DS_Store" || file.Name() == ".git" {
			continue
		}

		if isItLast(filePath, file, printFiles) {
			printGraph = dirGraph + "└───" + file.Name()
			tempGraph = "	"
		} else {
			printGraph = dirGraph + "├───"+file.Name()
			tempGraph = "│	"
		}

		if !file.IsDir(){
			if file.Size() == 0 {
				tempSize = " (empty)"
			} else {
				tempSize = fmt.Sprintf(" (%vb)",file.Size())
			}
			printGraph = printGraph + tempSize
		}

	 	_, err := fmt.Fprintln(out, printGraph)
		if err != nil {
			 fmt.Fprintf(os.Stderr, "Fprintln: %v\n", err)
		}

		if file.IsDir() {
			levelDir(out, filePath + "/" + file.Name(), printFiles, dirGraph + tempGraph)
		}
    }

	return err
}

//Ф-ия проверяет является файл или папка последней на конкретном уровне дерева
func isItLast (filePath string, f os.FileInfo, printFiles bool) bool {
	files, err := ioutil.ReadDir(filePath)
    if err != nil {
        log.Fatal(err)
    }

	var tempCountFiles int = 0

	for i, file := range files {
		if !file.IsDir() {
			tempCountFiles++
		} else {
			tempCountFiles = 0
		}

		if file.Name() == f.Name() && (i + 1) == len(files) {
			return true
		}
	}

	if (len(files)-1 - tempCountFiles) > 0 {
		if files[len(files)-1 - tempCountFiles].Name() == f.Name() && !printFiles {
			return true
		}
	}
	return false
}
