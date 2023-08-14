package main

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestEmpty(t *testing.T) {
	input := ``
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	_, err := parser.Read()
	if err != io.EOF {
		t.Fatal("EOF should be reached")
	}
}

func TestEmptyLine(t *testing.T) {
	input := `
	`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	_, err := parser.Read()
	if err != io.EOF {
		t.Fatal("EOF should be reached")
	}
}

func TestCommentLine(t *testing.T) {
	input := `#`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	_, err := parser.Read()
	if err != io.EOF {
		t.Fatal("EOF should be reached")
	}
}

func TestNoOwners(t *testing.T) {
	input := `file`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	_, err := parser.Read()
	parsingErr, isParsingError := err.(*ParsingErr)
	if !isParsingError {
		t.Fatal("Parsing error is expected")
	}
	if parsingErr.Line != 1 || parsingErr.Column != 4 {
		t.Error("Parsing error should be at 1:4")
	}
}

func TestSingleOwner(t *testing.T) {
	input := `file test@example.com`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	codeowner, err := parser.Read()
	if err != nil {
		t.Fatalf("No errors are expected %v", err)
	}
	if codeowner.Pattern != "file" {
		t.Error("Expected pattern is 'file'")
	}
	if len(codeowner.Owners) != 1 || codeowner.Owners[0] != "test@example.com" {
		t.Fatal("Expected owner is 'test@example.com'")
	}
}

func TestSingleOwnerWithComment(t *testing.T) {
	input := `file test@example.com # comment`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	codeowner, err := parser.Read()
	if err != nil {
		t.Fatalf("No errors are expected %v", err)
	}
	if codeowner.Pattern != "file" {
		t.Error("Expected pattern is 'file'")
	}
	if len(codeowner.Owners) != 1 {
		t.Fatal("1 codeowner is expected")
	}
	if codeowner.Owners[0] != "test@example.com" {
		t.Error("Expected owner is 'test@example.com'")
	}
}

func TestMultipleOwners(t *testing.T) {
	input := `file test@example.com other@example.com`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	codeowner, err := parser.Read()
	if err != nil {
		t.Fatalf("No errors are expected %v", err)
	}
	if codeowner.Pattern != "file" {
		t.Error("Expected pattern is 'file'")
	}
	if len(codeowner.Owners) != 2 {
		t.Fatal("2 owners are expected")
	}
	if codeowner.Owners[0] != "test@example.com" || codeowner.Owners[1] != "other@example.com" {
		t.Error("Expected owners are 'test@example.com' and 'other@example.com'")
	}
}

func TestInvalid(t *testing.T) {
	input := `invalid
	file test@example.com other@example.com`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	codeowner, err := parser.Read()
	_, isParsingError := err.(*ParsingErr)
	if !isParsingError {
		t.Fatal("Parsing error is expected")
	}
	codeowner, err = parser.Read()
	if err != nil {
		t.Fatalf("No errors are expected %v", err)
	}
	if codeowner.Pattern != "file" {
		t.Error("Expected pattern is 'file'")
	}
	if len(codeowner.Owners) != 2 {
		t.Fatal("2 owners are expected")
	}
	if codeowner.Owners[0] != "test@example.com" || codeowner.Owners[1] != "other@example.com" {
		t.Error("Expected owners are 'test@example.com' and 'other@example.com'")
	}
}

func TestInvalidReadAll(t *testing.T) {
	input := `invalid
	file test@example.com other@example.com`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	codeowners, parsingErrors, err := ReadAll(parser)
	if err != nil {
		t.Fatalf("No errors are expected %v", err)
	}
	if len(parsingErrors) != 1 {
		t.Fatal("1 parsing error is expected")
	}
	if len(codeowners) != 1 {
		t.Fatal("1 codeowners entry are expected")
	}
	if codeowners[0].Pattern != "file" {
		t.Error("Expected pattern is 'file'")
	}
	if len(codeowners[0].Owners) != 2 {
		t.Fatal("2 owners are expected")
	}
	if codeowners[0].Owners[0] != "test@example.com" || codeowners[0].Owners[1] != "other@example.com" {
		t.Error("Expected owners are 'test@example.com' and 'other@example.com'")
	}
}

func TestComplex(t *testing.T) {
	input := `file test@example.com other@example.com
	
	# comment line
	
	other test@example.com`
	parser := NewBufferedCodeownersParser(bufio.NewReader(bytes.NewReader([]byte(input))))
	codeowners, _, err := ReadAll(parser)
	if err != nil {
		t.Fatalf("No errors are expected %v", err)
	}
	if len(codeowners) != 2 {
		t.Fatal("2 codeowners entry are expected")
	}
	if codeowners[0].Pattern != "file" {
		t.Error("Expected pattern is 'file'")
	}
	if len(codeowners[0].Owners) != 2 {
		t.Fatal("2 owners are expected")
	}
	if codeowners[0].Owners[0] != "test@example.com" || codeowners[0].Owners[1] != "other@example.com" {
		t.Error("Expected owners are 'test@example.com' and 'other@example.com'")
	}
	if codeowners[1].Pattern != "other" {
		t.Error("Expected pattern is 'other'")
	}
	if len(codeowners[1].Owners) != 1 {
		t.Fatal("1 owner is expected")
	}
	if codeowners[1].Owners[0] != "test@example.com" {
		t.Error("Expected owners are 'test@example.com' and 'other@example.com'")
	}
}
