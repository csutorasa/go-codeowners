// Main application package.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var debugMode bool

var debugLogger = log.New(io.Discard, "", 0)

func main() {
	debug := flag.Bool("debug", false, "--debug to enable debug logging")
	gitDir := flag.String("C", "", "-C path/to/git/directory, defaults to current working directory")
	flag.Parse()
	if *debug {
		debugLogger = log.New(os.Stdout, "DEBUG ", log.Ldate|log.Ltime)
	}
	if *gitDir == "" {
		debugLogger.Print("Git directory was not found in arguments")
		workingDirectory, err := os.Getwd()
		if err != nil {
			panic("Failed to find current working directory")
		}
		debugLogger.Printf("Using current working directory %s to find git repository root", workingDirectory)
		*gitDir, err = FindRepositoryRoot(workingDirectory)
		if err != nil {
			debugLogger.Printf("Did not find git repository root, using current working directory %s", workingDirectory)
			*gitDir = workingDirectory
		}
	} else if debugMode {
		debugLogger.Printf("Using git directory from arguments %s", *gitDir)
	}
	filePaths := flag.Args()
	if len(filePaths) == 0 {
		panic("Missing file arguments")
	}
	err := findCodeowners(*gitDir, filePaths)
	if err != nil {
		panic(err)
	}
}

func findCodeowners(gitDir string, filePaths []string) error {
	codeownersFilepath := FindCodeownersFile(gitDir)
	if codeownersFilepath == "" {
		return fmt.Errorf("failed to find CODEOWNERS in %s", gitDir)
	}
	debugLogger.Printf("Found CODEOWNERS file %s", codeownersFilepath)
	codeowners, err := ValidateCodeownersFile(codeownersFilepath)
	if err != nil {
		return err
	}
	for _, filePath := range filePaths {
		err := doCheck(codeowners, gitDir, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func doCheck(codeowners []Codeowner, gitDir string, filePath string) error {
	exists, isDir := FindFile(gitDir, filePath)
	if !exists {
		fmt.Fprintf(os.Stderr, "Cannot find %s, will consider it an existing file\n", filePath)
	}
	owners := FindOwners(codeowners, filePath, isDir)
	if owners == nil {
		fmt.Printf("%s has no codeowners entry\n", filePath)
	} else {
		fmt.Printf("%s is owned by %v\n", filePath, owners)
	}
	return nil
}
