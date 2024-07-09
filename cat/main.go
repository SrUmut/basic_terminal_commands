package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	argc, argv := len(os.Args), os.Args

	files := make([]*os.File, 0)
	displayLineNumbers := false

	if argc < 1 {
		fmt.Println("Error: No file to cat.")
		os.Exit(1)
	}
	for _, arg := range argv[1:] {
		//if argument is a flag
		if strings.HasPrefix(arg, "-") {
			for _, flag := range arg[1:] {
				switch flag {
				case 'n':
					displayLineNumbers = true
				default:
					fmt.Printf("Error: no flag as %v\n", flag)
				}
			}
		} else { // if argument is a path to a file😍
			file, err := os.Open(arg)
			checkError(err)
			defer file.Close()
			files = append(files, file)
		}
	}

	readBuffer := make([]byte, 1)
	lineCounter := 1
	prevLineFeed := true
	for _, file := range files {
		// read file until EOF
		for {
			_, err := file.Read(readBuffer)
			if err == io.EOF {
				break
			}
			checkError(err)
			if displayLineNumbers { // -n flag
				if prevLineFeed {
					fmt.Printf("%6v  %s", lineCounter, readBuffer)
					lineCounter++
					prevLineFeed = false
				} else {
					fmt.Printf("%s", readBuffer)
				}
				if readBuffer[0] == '\n' {
					prevLineFeed = true
				}
			} else {
				fmt.Printf("%s", readBuffer)
			}
		}
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
