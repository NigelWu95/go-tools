package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func help() {
	fmt.Printf("Usage: qsresult <source-dir-path> <target-dir-path>\r\n")
}

func main() {

	args := os.Args
	argc := len(args)

	var sourcePath string
	var targetPath string

	switch argc {
	case 3:
		sourcePath = args[1]
		targetPath = args[2]
	default:
		help()
		return
	}

	resultFilePath := string(sourcePath + string(filepath.Separator) + "result.txt")
	resultFile, err := os.Open(resultFilePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer resultFile.Close()

	var successLine string
	var order string
	var mkdirErr error
	var renameErr error
	br := bufio.NewReader(resultFile)
	for {
		line, _, e := br.ReadLine()
		if e == io.EOF {
			break
		}
		fmt.Println(string(line))
		successLine = strings.Split(string(line), ": ")[0]
		order = strings.Split(successLine, " ")[1]
		successFilePath := sourcePath + string(filepath.Separator) + "listbucket_success_" + order + ".txt"
		mkdirErr = os.Mkdir(targetPath, 0)
		if mkdirErr != nil {
			fmt.Printf("mkdir %s error: %s\n", targetPath, err)
		}
		renameErr = os.Rename(successFilePath, targetPath + string(filepath.Separator) + order + ".txt")
		if renameErr != nil {
			fmt.Printf("move %s to %s error: %s\n", successFilePath, targetPath, err)
		}
	}
}
