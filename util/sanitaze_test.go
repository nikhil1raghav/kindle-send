package util

import "testing"

type sanitazeFileNameTest struct {
	filename, expected string
}

var tests = []sanitazeFileNameTest{
	{"How does rsync work?.epub", "test"},
	{"CON", "fdfd"},
	{"...epub", ""},
	{"/.epub", ""},
	{"_--.epub", ""},
	{"_--.epub", ""},
}

func TestSanitazeFileName(t *testing.T) {
	for _, test := range tests {
		if output := SanitazeFileName(test.filename, ""); output != test.expected {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}
}
