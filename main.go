package main

import (
	"fmt"
	"io"
	"os"
	// "path/filepath"
	// "strings
	"io/ioutil"
	"log"
)

func printing(filePath string, f os.FileInfo, printFiles bool, levelInd int, lastInd int){

	var tempGraph string
	var size string
	var tempi int
	// var tempi int

	if !(isItLast(filePath, f, printFiles)) && (lastInd > 0) {
		tempi = levelInd
		levelInd = levelInd - lastInd
		lastInd = tempi
	}

	if isItLast(filePath, f, printFiles) && (lastInd > 0) {
		tempi = levelInd
		levelInd = levelInd - (lastInd - 1)
		lastInd = tempi - levelInd + 1
	}

	for i := 0; i < levelInd; i++ {
		tempGraph = tempGraph + "│   "
	}

	for j := 0; j < lastInd-1; j++ {
		tempGraph = tempGraph + "    "
	}

	if isItLast(filePath, f, printFiles) {
		tempGraph = tempGraph + "└───" + f.Name()
	} else {
		tempGraph = tempGraph + "├───"+f.Name()
	}

	if !f.IsDir(){
		if f.Size() == 0 {
			size = " (empty)"
		} else {
			size = fmt.Sprintf(" (%vb)",f.Size())
		}
		tempGraph = tempGraph + size
	}

	 fmt.Println(tempGraph)
}

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

func levelDir(filePath string, printFiles bool, levelInd int, lastInd int){
	files, err := ioutil.ReadDir(filePath)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
		if !file.IsDir() && !printFiles {
			continue
		}
		if isItLast(filePath, file, printFiles) {
			lastInd = lastInd + 1
		}

		printing(filePath, file, printFiles, levelInd, lastInd)

		if file.IsDir() {
			levelDir(filePath + "/" + file.Name(), printFiles, levelInd+1, lastInd)
		}
    }
}

func dirTree(out io.Writer, filePath string, printFiles bool) error  {
	var levelInd = 0
	var lastInd = 0
	levelDir(filePath, printFiles, levelInd, lastInd)

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
