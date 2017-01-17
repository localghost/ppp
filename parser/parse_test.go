package parser

import (
	"encoding/json"
	"os"
	"testing"
	"reflect"
	"path/filepath"
)

func decode(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result interface{}
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

var testCases = []struct {
	entry string
	expected string
}{
	{"001-main-single-inline.json", "001-expected.json"},
	{"002-main-array-inline.json", "002-expected.json"},
	{"003-main-nested-inline.json", "003-expected.json"},
	{"004-main-subdir-inline.json", "004-expected.json"},
	{"005-main-action-shell-single.json", "005-expected.json"},
	{"006-main-action-shell-array.json", "006-expected.json"},
}

func TestRunner(t *testing.T) {
	testRunner := func (entryPath string, expectedPath string) {
		fixtureDir := "./test-fixtures"

		actual, err := NewParser().ParseFile(filepath.Join(fixtureDir, entryPath))
		if err != nil {
			t.Fatal(err)
		}

		expected, err := decode(filepath.Join(fixtureDir, expectedPath))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatal(expected, actual)
		}
	}

	for _, tc := range testCases {
		testRunner(tc.entry, tc.expected)
	}
}