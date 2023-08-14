package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var knownLocations = []string{
	"CODEOWNERS",
	"docs/CODEOWNERS",
	".github/CODEOWNERS",
}

const fileSizeLimit int64 = 3 * 1024 * 1024

func FindRepositoryRoot(d string) (string, error) {
	command := exec.Command("git", "rev-parse", "--show-toplevel")
	err := command.Run()
	if err != nil {
		return "", err
	}
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func FindCodeownersFile(d string) string {
	for _, knownLocation := range knownLocations {
		p := filepath.Join(d, knownLocation)
		_, err := os.Stat(p)
		if err != nil {
			continue
		}
		return p
	}
	return ""
}

func ValidateCodeownersFile(f string) ([]Codeowner, error) {
	fileInfo, err := os.Stat(f)
	if err != nil {
		return nil, err
	}
	if fileInfo.Size() > fileSizeLimit {
		fmt.Fprint(os.Stderr, "CODEOWNERS is larger than 3 MB")
	}
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	parser := NewBufferedCodeownersParser(bufio.NewReader(file))
	codeowners, parsingErrors, err := ReadAll(parser)
	if err != nil {
		return nil, err
	}
	for _, parsingError := range parsingErrors {
		fmt.Fprintf(os.Stderr, "%v\n", parsingError.Error())
	}
	return codeowners, nil
}

func FindFile(d string, f string) (bool, bool) {
	fp := filepath.Join(d, f)
	fileInfo, err := os.Stat(fp)
	if err != nil {
		return false, false
	}
	return true, fileInfo.IsDir()
}
