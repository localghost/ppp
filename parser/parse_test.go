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

type testCase struct {
	tested string
	expected string
}

func runTest(t *testing.T, tc *testCase) {
	fixtureDir := "./test-fixtures"

	actual, err := NewParser().ParseFile(filepath.Join(fixtureDir, tc.tested))
	if err != nil {
		t.Fatal(err)
	}

	expected, err := decode(filepath.Join(fixtureDir, tc.expected))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal(expected, actual)
	}
}

func runTests(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		runTest(t, &tc)
	}
}

func TestInline(t *testing.T) {
	testCases := []testCase {
		{"001-main-single-inline.json", "001-expected.json"},
		{"002-main-array-inline.json", "002-expected.json"},
		{"003-main-nested-inline.json", "003-expected.json"},
		{"004-main-subdir-inline.json", "004-expected.json"},
	}
	runTests(t, testCases)
}

func TestShell(t *testing.T) {
	testCases := []testCase {
		{"005-main-action-shell-single.json", "005-expected.json"},
		{"006-main-action-shell-array.json", "006-expected.json"},
	}
	runTests(t, testCases)
}

func TestExec(t *testing.T) {
	testCases := []testCase {
		{"007-main-action-exec-single.json", "007-expected.json"},
	}
	runTests(t, testCases)
}