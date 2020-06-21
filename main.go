// B''H

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var gitDirs []string

func main() {

	var startDirs []string = make([]string, 1)
	var err error

	// Determine what the starting directory will be. If not specified, it will be the current directory
	if len(os.Args) == 1 {

		startDirs[0], err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		startDirs = os.Args[1:]
	}

	// Search for all the directories that have a .git folder in them
	for _, startDir := range startDirs {
		startDir, err = filepath.Abs(startDir)
		if err != nil {
			log.Fatal(err)
		}

		search(startDir)
	}

	// Say how many repos found
	fmt.Printf("%d repositories found\n", len(gitDirs))

	// Go into each directory and run git status
	for _, gitDir := range gitDirs {
		doGitStatus(gitDir)
	}
}

func search(directory string) {

	// Get the list of files in the search directory
	files := fetchFiles(directory)
	var dirBool bool

	for _, file := range files {

		var fullPath string = directory + "/" + file.Name()

		dirBool = dirTest(fullPath)

		if ".git" == strings.ToLower(file.Name()) && dirBool {
			gitDirs = append(gitDirs, directory)
		}

		if dirBool {
			search(fullPath)
		}

	}

}

func fetchFiles(directory string) []os.FileInfo {
	var files []os.FileInfo
	var err error

	files, err = ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func dirTest(filename string) bool {

	fileStat, err := os.Lstat(filename)
	if err != nil {
		log.Fatal(err)
	}

	fileMode := fileStat.Mode()

	return fileMode.IsDir()

}

func doGitStatus(directory string) {

	var err error
	var cmd *exec.Cmd
	var stdout []byte

	err = os.Chdir(directory)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("~==========~")
	fmt.Println(directory)
	fmt.Println("~==========~")

	cmd = exec.Command("git", "status")

	stdout, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(stdout))

}
