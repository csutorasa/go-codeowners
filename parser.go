package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ParsingErr struct {
	Line    int
	Column  int
	Message string
}

func (e *ParsingErr) Error() string {
	return fmt.Sprintf("parsing error found at %d:%d - %s", e.Line, e.Column, e.Message)
}

func ReadAll(p CodeownersParser) ([]Codeowner, []*ParsingErr, error) {
	codeowners := []Codeowner{}
	errors := []*ParsingErr{}
	for {
		codeowner, err := p.Read()
		if err != nil {
			parsingErr, isParsingError := err.(*ParsingErr)
			if isParsingError {
				errors = append(errors, parsingErr)
				continue
			}
			if err == io.EOF {
				return codeowners, errors, nil
			}
			return codeowners, errors, err
		}
		codeowners = append(codeowners, codeowner)
	}
}

type CodeownersParser interface {
	Read() (Codeowner, error)
}

type bufferedCodeownersParser struct {
	buf       *bufio.Reader
	lineCount int
}

func NewBufferedCodeownersParser(buf *bufio.Reader) CodeownersParser {
	return &bufferedCodeownersParser{
		buf:       buf,
		lineCount: 0,
	}
}

func (p *bufferedCodeownersParser) Read() (Codeowner, error) {
	var codeowner Codeowner
	input, err := p.readUntilNonEmptyLine()
	if err != nil {
		return codeowner, err
	}
	parts := strings.Split(input, " ")
	if len(parts) == 1 {
		return codeowner, &ParsingErr{
			Line:    p.lineCount,
			Column:  len(input),
			Message: "No space was found",
		}
	}
	codeowner.Pattern = Pattern(parts[0])
	codeowner.Owners = make([]Owner, len(parts)-1)
	for i, owner := range parts[1:] {
		codeowner.Owners[i] = Owner(owner)
	}
	return codeowner, nil
}

func (p *bufferedCodeownersParser) readUntilNonEmptyLine() (string, error) {
	for {
		input, err := p.buf.ReadString('\n')
		p.lineCount++
		if (err != nil && err != io.EOF) || (err == io.EOF && len(input) == 0) {
			return "", err
		}
		commentStarts := strings.Index(input, "#")
		if commentStarts >= 0 {
			input = input[:commentStarts]
		}
		input = strings.TrimSpace(input)
		if len(input) > 0 {
			return input, nil
		}
	}
}
