package parser

import (
	"encoding/json"
	"os"
	"testing"
	"reflect"
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

func TestInline(t *testing.T) {
	actual, err := ParseFile("./test-fixtures/01-main.json")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := decode("./test-fixtures/01-expected.json")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatal(expected, actual)
	}
}