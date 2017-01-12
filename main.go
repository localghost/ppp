package main

import (
	"fmt"
	"github.com/localghost/packer-pp/parser"
	"os"
	"encoding/json"
)

func main() {
	templatePath := os.Args[1]
	result, err := parser.ParseFile(templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
	json.NewEncoder(os.Stdout).Encode(result)
}
