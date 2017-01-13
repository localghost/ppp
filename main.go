package main

import (
	"fmt"
	"github.com/localghost/ppp/parser"
	"os"
	"encoding/json"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: ppp <template-path>")
		os.Exit(1)
	}

	templatePath := os.Args[1]
	result, err := parser.ParseFile(templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(result)
}
