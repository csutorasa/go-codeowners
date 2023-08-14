package main

import (
	"path/filepath"
	"strings"
)

const DoubleAsterisk = "**"
const Asterisk = "*"
const Wildcard = "?"

func FindOwners(codeowners []Codeowner, f string, isDir bool) []Owner {
	for i := len(codeowners) - 1; i >= 0; i-- {
		codeowner := codeowners[i]
		if (!codeowner.Pattern.ContainsNonTrailingSlash() || !isDir) && codeowner.Pattern.Matches(f) {
			return codeowner.Owners
		}
	}
	return nil
}

func (p Pattern) Matches(fp string) bool {
	if p.ContainsNonTrailingSlash() {
		return matches(p.splitToParts(), fp)
	} else {
		fileName := filepath.Base(fp)
		return matches(p.splitToParts(), fileName)
	}
}

func (p Pattern) ContainsNonTrailingSlash() bool {
	s := string(p)
	if s == "" {
		return false
	}
	index := strings.Index(s, "/")
	return index != -1 && index != len(s)-1
}

func (p Pattern) RemoveTrailingSlash() string {
	s := string(p)
	if s == "" {
		return s
	}
	if s[len(s)-1] == '/' {
		return s[:len(s)-1]
	}
	return s
}

func (p Pattern) splitToParts() []string {
	s := p.RemoveTrailingSlash()
	parts := []string{}
	doubleAsterisks := strings.Split(s, DoubleAsterisk)
	for i, doubleAsterisk := range doubleAsterisks {
		asterisks := strings.Split(doubleAsterisk, Asterisk)
		for j, asterisk := range asterisks {
			wildcards := strings.Split(asterisk, Wildcard)
			for k, wildcard := range wildcards {
				if wildcard != "" {
					parts = append(parts, wildcard)
				}
				if k != len(wildcards)-1 {
					parts = append(parts, Wildcard)
				}
			}
			if j != len(asterisks)-1 {
				parts = append(parts, Asterisk)
			}
		}
		if i != len(doubleAsterisks)-1 {
			parts = append(parts, DoubleAsterisk)
		}
	}
	return parts
}

func matches(parts []string, f string) bool {
	index := 0
	if len(parts) == 0 {
		return true
	}
	for i, part := range parts {
		if part == DoubleAsterisk {
			if len(parts)-1 == i {
				return true
			}
			for j := 0; j < len(f)-index; j++ {
				if matches(parts[i+1:], f[index+j:]) {
					return true
				}
			}
		} else if part == Asterisk {
			if len(parts)-1 == i && !strings.Contains(f[index:], "/") {
				return true
			}
			for j := 0; j < len(f)-index && f[j+index] != '/'; j++ {
				if matches(parts[i+1:], f[index+j:]) {
					return true
				}
			}
		} else if part == Wildcard {
			if index+1 > len(f) {
				return false
			}
			index++
		} else {
			if index+len(part) > len(f) {
				return false
			}
			if f[index:index+len(part)] != part {
				return false
			}
			index += len(part)
		}
	}
	return index == len(f)
}
