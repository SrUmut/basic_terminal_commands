package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func colored(text, color string) string {
	switch color {
	case "red":
		return fmt.Sprintf("\x1b[31m%v\x1b[0m", text)
	case "cyan":
		return fmt.Sprintf("\033[36m%v\033[0m", text)
	default:
		fmt.Printf("Error: No such color as %v\n", color)
		return ""
	}
}

func getHumanReadable(byteSize int64) string {
	var hr string
	byteSizeF := float32(byteSize)
	if byteSizeF >= 1125899906842624 { // PB
		hr = fmt.Sprintf("%.1fP", byteSizeF/1099511627776)
	} else if byteSizeF >= 1099511627776 { // TB
		hr = fmt.Sprintf("%.1fT", byteSizeF/1099511627776)
	} else if byteSizeF >= 1073741824 { // GB
		hr = fmt.Sprintf("%.1fG", byteSizeF/1073741824)
	} else if byteSizeF >= 1048576 { // MB
		hr = fmt.Sprintf("%.1fM", byteSizeF/1048576)
	} else if byteSizeF >= 1024 { // KB
		hr = fmt.Sprintf("%.1fK", byteSizeF/1024)
	} else {
		hr = fmt.Sprintf("%vB", byteSize)
	}
	return hr
}

func main() {
	argc, argv := len(os.Args), os.Args
	path := "./"
	longFormat := false // -l
	listAll := false    // -a
	humanReadable := false
	if argc > 1 {
		for _, arg := range argv[1:] {
			if strings.HasPrefix(arg, "-") { // if argument is a flag, exp: -l
				for _, flag := range arg[1:] {
					switch flag {
					case 'l':
						longFormat = true
					case 'a':
						listAll = true
					case 'h':
						humanReadable = true
					default:
						fmt.Printf("Error: no flag as %v\n", flag)
						os.Exit(1)
					}
				}
			} else { // if argument is a dir path check if directory exists
				fileinfo, err := os.Stat(arg)
				checkError(err)
				if !fileinfo.IsDir() {
					checkError(fmt.Errorf("given path does not belong to a directory: %v", path))
				}
				path = arg
			}
		}
	}
	// get the list of entries in the directory
	list, err := os.ReadDir(path)
	checkError(err)
	for _, entry := range list {
		// -a
		if !listAll && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		filePath := path + "/" + entry.Name()
		fileInfo, err := os.Stat(filePath)
		checkError(err)

		// getting entry permissions as string
		entryPerms := ""
		if entry.IsDir() {
			entryPerms = "d" + fileInfo.Mode().Perm().String()[1:]
		} else {
			entryPerms = fileInfo.Mode().Perm().String()
		}

		// -l
		if longFormat {
			fmt.Printf("%v\t", entryPerms)

			// Entry Size
			entrySize := fileInfo.Size()
			if humanReadable {
				fmt.Printf("%6v\t", getHumanReadable(entrySize))
			} else {
				fmt.Printf("%v\t", entrySize)
			}

			// Modification Time
			modTime := fileInfo.ModTime()
			fmt.Printf("%v ", modTime.Format("Jan"))
			fmt.Printf("%2v ", modTime.Format("2"))
			if time.Now().Year()-modTime.Year() >= 1 {
				fmt.Printf("%5v\t", modTime.Format("2006"))
			} else {
				fmt.Printf("%5v\t", modTime.Format("15:04"))
			}
		}

		// Setting entry name
		entryName := ""
		if entry.IsDir() { // directory
			entryName = colored(entry.Name(), "cyan")
		} else { // file
			if strings.Contains(entryPerms, "x") { // executable
				entryName = colored(entry.Name(), "red")
			} else { // not executable
				entryName = entry.Name()
			}
		}
		fmt.Printf("%v\n", entryName)
	}
}
