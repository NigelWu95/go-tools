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

	//args := os.Args
	//args = append(args, "../temp2")
	//args = append(args, "../temp3")
	//argc := len(args)

	var sourcePath string
	var targetPath string

	switch argc {
	case 3:
		sourcePath = args[1]
		targetPath = args[2]
		exist, err := PathExists(targetPath)
		if err == nil {
			if !exist {
				mkdirErr := os.Mkdir(targetPath, 0755)
				if mkdirErr != nil {
					fmt.Printf("mkdir %s error: %s\n", targetPath, err.Error())
				}
			}
		} else {
			fmt.Printf("Error: %s\n", err.Error())
		}
		//syscall.Umask(0)
		//err = os.Chmod(targetPath, 0755)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	default:
		help()
		return
	}

	resultFilePath := string(sourcePath + string(filepath.Separator) + "result.txt")
	exist, err := PathExists(resultFilePath)
	if !exist {
		fmt.Printf("no more finished files, error: %s\n", err.Error())
		return
	}
	resultFile, err := os.Open(resultFilePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	var lineItems []string
	var order string
	var successFilePath string
	var rewriteLines []string
	var successMap = map[string]string{}
	br := bufio.NewReader(resultFile)
	for {
		line, _, e := br.ReadLine()
		if e == io.EOF {
			break
		}

		if strings.Contains(string(line), "successfully done") {
			lineItems = strings.Split(string(line), " ")
			order = strings.Split(lineItems[1], ":")[0]
			successFilePath = sourcePath + string(filepath.Separator) + "listbucket_success_" + order + ".txt"
			err = os.Rename(successFilePath, targetPath + string(filepath.Separator) + order + ".txt")
			if err != nil {
				rewriteLines = append(rewriteLines, string(line))
				fmt.Printf("move %s to %s error: %s\n", successFilePath, targetPath, err.Error())
			} else {
				successMap[strings.Join(lineItems[0:2], " ")] = strings.Split(lineItems[2], "\t")[0]
			}
		} else {
			rewriteLines = append(rewriteLines, string(line))
		}
	}

	closeErr := resultFile.Close()
	if closeErr != nil {
		fmt.Printf("close \"result.txt\" error: %s\n", err.Error())
	}
	resultFile, err = os.Create(resultFilePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	writer := bufio.NewWriter(resultFile)

	var key string
	for _,line := range rewriteLines {
		lineItems = strings.Split(line, " ")
		key = strings.Join(lineItems[0:2], " ")
		if _, ok := successMap[key]; !ok {
			_, err = fmt.Fprintln(writer, line)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return
			}
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
}
