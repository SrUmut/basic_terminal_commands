package main

import (
	"fmt"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	argc, argv := len(os.Args), os.Args
	path := ""
	if argc > 1 {
		path = argv[1]
	} else {
		path = "./"
	}
	list, err := os.ReadDir(path)
	checkError(err)
	for _, entry := range list {
		/*if entry.IsDir() {
			fmt.Println("\033[36m" + entry.Name() + "\033[0m")
		} else {
			//fmt.Println(entry.Name())
		}*/
		filePath := path + "/" + entry.Name()
		fileInfo, err := os.Stat(filePath)
		checkError(err)
		fmt.Println(fileInfo.Mode().Perm(), entry.Name())
	}
}
