package main

import (
	"testing"
)

func TestFindOwners(t *testing.T) {
	owners := FindOwners([]Codeowner{
		{
			Pattern: Pattern("test.go"),
			Owners: []Owner{
				Owner("test@example.com"),
			},
		},
	}, "test.go", false)
	if len(owners) != 1 {
		t.Fatal("1 owner is expected")
	}
	if string(owners[0]) != "test@example.com" {
		t.Fatal("Expected owner is 'test@example.com'")
	}
}

func TestExactMatch(t *testing.T) {
	pattern := Pattern("file.go")
	matches := pattern.Matches("file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("file2.go")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file2.go")
	if matches {
		t.Error("Exact match is not expected")
	}
}

func TestWildcardMatch(t *testing.T) {
	pattern := Pattern("fil?.go")
	matches := pattern.Matches("file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("file2.go")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file2.go")
	if matches {
		t.Error("Exact match is not expected")
	}
}

func TestTrailingSlashMatch(t *testing.T) {
	pattern := Pattern("dir/")
	matches := pattern.Matches("dir")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("dir/dir")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/dir2")
	if matches {
		t.Error("Exact match is not expected")
	}
}

func TestAsteriskMatch(t *testing.T) {
	pattern := Pattern("*.go")
	matches := pattern.Matches("file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
}

func TestAsteriskEndingMatch(t *testing.T) {
	pattern := Pattern("file.*")
	matches := pattern.Matches("file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("file2.go")
	if matches {
		t.Error("Exact is not expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file2.go")
	if matches {
		t.Error("Exact is not expected")
	}
}

func TestMultipleAsteriskMatch(t *testing.T) {
	pattern := Pattern("*.*")
	matches := pattern.Matches("file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("file.go2")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file.go2")
	if !matches {
		t.Error("Exact match expected")
	}
}

func TestMultiLevelMatch(t *testing.T) {
	pattern := Pattern("dir/file.go")
	matches := pattern.Matches("file.go")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("test/dir/file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
}

func TestMultiLevelWildcardMatch(t *testing.T) {
	pattern := Pattern("dir/fil?.go")
	matches := pattern.Matches("file.go")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("test/dir/file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
}

func TestMultiLevelAsterixMatch(t *testing.T) {
	pattern := Pattern("dir/*.go")
	matches := pattern.Matches("file.go")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("test/dir/file.go")
	if matches {
		t.Error("Exact match is not expected")
	}
}

func TestMultiLevelDoubleAsterixMatch(t *testing.T) {
	pattern := Pattern("**/*.go")
	matches := pattern.Matches("file.go")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("dir/file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
	matches = pattern.Matches("test/dir/file.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("test/dir/file.go2")
	if matches {
		t.Error("Exact match is not expected")
	}
}

func TestWildcardRequiredMatch(t *testing.T) {
	pattern := Pattern("file.go?")
	matches := pattern.Matches("file.go2")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("file.go")
	if matches {
		t.Error("Exact match is not expected")
	}
	pattern = Pattern("file*?*.go")
	matches = pattern.Matches("file2.go")
	if !matches {
		t.Error("Exact match expected")
	}
	matches = pattern.Matches("file.go")
	if matches {
		t.Error("Exact match is not expected")
	}
}
