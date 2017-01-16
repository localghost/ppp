package parser

import "fmt"

// Appends \n to format and forwards call to fmt.Printf().
func println(format string, a ...interface{}) (int, error) {
	return fmt.Printf(format + "\n", a)
}

// Copies entries from src map which are not present in dest map.
func merge(dest map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		if _, ok := dest[k]; !ok {
			dest[k] = v
		}
	}
}

func getIndexOr(slice interface{}, index int, defaultValue interface{}) interface{} {
	switch slice := slice.(type) {
	case []string:
		if len(slice) > index {
			return slice[index]
		}
	}
	return defaultValue
}
