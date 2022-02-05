package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func outDir(path string, d fs.DirEntry, str []string, gap int, out io.Writer) (int, []string) {
	files, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		panic(err)
	}
	for ind, file := range files {
		//последний файл
		if file.Name() == d.Name() && ind == len(files)-1 {
			//если дирикторий
			if d.IsDir() {

				str[len(str)-1] = "└───"

				for _, elem := range str {
					fmt.Fprint(out, elem)
				}

				gap += 1
				fmt.Fprintln(out, d.Name())
				str[len(str)-1] = "	"
				str = append(str, "├───")
				return gap, str
			}
			//если не дирикторий
			str[len(str)-1] = "└───"

			for _, elem := range str {
				fmt.Fprint(out, elem)
			}

			if val := fileSize(path); val != 0 {
				fmt.Fprintln(out, fmt.Sprintf("%s (%vb)", d.Name(), fileSize(path)))
			} else {
				fmt.Fprintln(out, fmt.Sprintf("%s (empty)", d.Name()))
			}

			if gap+1 >= len(str) {
				return gap, str
			}
			for i := 0; i < gap+1; i++ {
				str = str[:len(str)-1]

			}
			str[len(str)-1] = "├───"
			return gap, str
		}
	}
	// если дирикторий не последний
	if d.IsDir() {

		for _, elem := range str {
			fmt.Fprint(out, elem)
		}

		fmt.Fprintln(out, d.Name())
		str[len(str)-1] = "│	"
		str = append(str, "├───")
		gap = 0

		//если документ
	} else {
		for _, elem := range str {
			fmt.Fprint(out, elem)
		}

		if val := fileSize(path); val != 0 {
			fmt.Fprintln(out, fmt.Sprintf("%s (%vb)", d.Name(), fileSize(path)))
		} else {
			fmt.Fprintln(out, fmt.Sprintf("%s (empty)", d.Name()))
		}
	}
	return gap, str
}

func fileSize(path string) int64 {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	fi, err := file.Stat()

	if err != nil {
		panic(err)
	}

	return fi.Size()
}

func CountDirs(path string) int {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	numOfDirs := 0
	for _, file := range files {
		if file.IsDir() {
			numOfDirs++
		}
	}
	return numOfDirs
}

func outDirFile(path string, d fs.DirEntry, str []string, gap int, res int, out io.Writer) (int, []string) {
	//fmt.Println(path)
	files, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		panic(err)
	}
	numOfDirs := CountDirs(filepath.Dir(path))

	count := 0
	for _, file := range files {
		//учет номера дириктория
		if file.IsDir() {
			count++
		}
		//если последний дирикторий
		if file.Name() == d.Name() && count == numOfDirs && d.IsDir() {
			//если пустой последний
			if CountDirs(path) == 0 {
				str[len(str)-1] = "└───"

				for _, elem := range str {
					fmt.Fprint(out, elem)
				}

				fmt.Fprintln(out, d.Name())
				if gap+1 >= len(str) {
					return gap, str
				}
				for i := 0; i < gap+1; i++ {
					str = str[:len(str)-1]
				}
				str[len(str)-1] = "├───"
				return gap, str
			}
			//если непустой последний
			str[len(str)-1] = "└───"

			for _, elem := range str {
				fmt.Fprint(out, elem)
			}
			fmt.Fprintln(out, d.Name())
			str[len(str)-1] = "	"
			str = append(str, "├───")
			gap++
			return gap, str

			//если непоследний
		}
	}
	if d.IsDir() {
		//пустой
		if CountDirs(path) == 0 {
			for _, elem := range str {
				fmt.Fprint(out, elem)
			}
			fmt.Fprintln(out, d.Name())
			//непустой
		} else {
			for _, elem := range str {
				fmt.Fprint(out, elem)
			}
			fmt.Fprintln(out, d.Name())
			str[len(str)-1] = "│	"
			str = append(str, "├───")
			gap = 0
			return gap, str
		}
	}
	return gap, str
}

func dirTree(out io.Writer, filePath string, printFiles bool) error {
	//rempath := path
	fileSystem := os.DirFS(".")
	index := 0
	tempindex := 0
	gap := 0
	str := []string{"├───"}
	res := 0
	mapa := map[string]bool{
		"testdata":true,
		".":true,
		"main_test.go":true,
		"main.go":true,
		"hw1.md":true,
		"dockerfile":true,
	}
	fs.WalkDir(fileSystem, filePath, func(path string, d fs.DirEntry, err error) error {
		if printFiles {
			
			if _, inMap := mapa[filepath.Base(path)];inMap{
				return nil
			}
		} else { 
			if filepath.Base(path) == "." || filepath.Base(path) == "testdata" {
				return nil
			}
		}

		if err != nil {
			panic(err)
		}

		index = strings.Count(path, string(os.PathSeparator))
		res = index - tempindex
		if index < tempindex {
			gap = 0
		}

		tempindex = index
		if printFiles {
			gap, str = outDir(path, d, str, gap, out)
		} else {

			gap, str = outDirFile(path, d, str, gap, res, out)
		}
		return nil
	})

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
