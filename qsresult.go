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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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
		exist, err := PathExists(targetPath)
		if err == nil {
			if !exist {
				mkdirErr := os.Mkdir(targetPath, 0)
				if mkdirErr != nil {
					fmt.Printf("mkdir %s error: %s\n", targetPath, err)
				}
			}
		} else {
			fmt.Printf("Error: %s\n", err)
		}

	default:
		help()
		return
	}

	resultFilePath := string(sourcePath + string(filepath.Separator) + "result.txt")
	exist, err := PathExists(resultFilePath)
	if !exist {
		fmt.Printf("no more finished files.")
		return
	}
	resultFile, err := os.Open(resultFilePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	var strLine string
	var successLine string
	var order string
	var successFilePath string
	var renameErr error
	var unmoved []string
	br := bufio.NewReader(resultFile)
	for {
		line, _, e := br.ReadLine()
		if e == io.EOF {
			break
		}
		strLine = string(line)
		if strings.Contains(strLine, "successfully done") {
			successLine = strings.Split(string(line), ": ")[0]
			order = strings.Split(successLine, " ")[1]
			successFilePath = sourcePath + string(filepath.Separator) + "listbucket_success_" + order + ".txt"
			renameErr = os.Rename(successFilePath, targetPath+string(filepath.Separator)+order+".txt")
			if renameErr != nil {
				unmoved = append(unmoved, string(line))
				fmt.Printf("move %s to %s error: %s\n", successFilePath, targetPath, err)
			}
		}
	}
	closeErr := resultFile.Close()
	if closeErr != nil {
		fmt.Printf("close \"result.txt\" error: %s\n", err)
	}

	if len(unmoved) > 0 {
		resultFile, err = os.Create(resultFilePath)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		_, err = resultFile.WriteString(strings.Join(unmoved, "\n"))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	} else {
		err = os.Remove(resultFilePath)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}
}
