package regex

import (
	"testing"
)

func Test_trimRegex(t *testing.T) {
	if trimRegex("a**b*cf**e***d*") != "a*b*cf*e*d*" {
		t.Fail()
	}
}

func Test_regexMatch(t *testing.T) {
	var tests = []struct {
		toMatch string
		regex   string
		answer  bool
	}{
		{toMatch: "aab", regex: "a*b", answer: true},
		{toMatch: "aaaaab", regex: "a*b", answer: true},
		{toMatch: "b", regex: "a*b", answer: true},
		{toMatch: "aabbcdeee", regex: "aa*b*cdee*", answer: true},
		{toMatch: "aabc", regex: "a*b", answer: false},
		{toMatch: "aabce", regex: "a*b", answer: false},
	}
	for i, test := range tests {
		if test.answer != regexMatch(test.toMatch, test.regex) {
			t.Fatalf("Test number %d failed", i+1)
		}
	}
}
