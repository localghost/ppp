package parser

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type context struct {
	templatePath string
}

func merge(dest map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		if _, ok := dest[k]; !ok {
			dest[k] = v
		}
	}
}

// TODO: Verify there won't be problems with too deep recursion.
func walk(ctx *context, raw interface{}) error {
	switch node := raw.(type) {
	case []interface{}:
		for _, v := range node {
			if err := walk(ctx, v); err != nil {
				return err
			}
		}
	case map[string]interface{}:
		for k, v := range node {
			if transformation, ok := transformations[k]; ok {
				delete(node, k)

				data, err := transformation(ctx, v)
				if err != nil {
					return err
				}
				if data != nil {
					merge(node, data.(map[string]interface{}))
				}

				// Transformation changes the template, walk must be restarted.
				return walk(ctx, raw)
			}

			if err := walk(ctx, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func parse(ctx *context, r io.Reader) (interface{}, error) {
	var result interface{}
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		return nil, err
	}

	if err := walk(ctx, result); err != nil {
		return nil, err
	}

	return result, nil
}

func ParseFile(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}

	return parse(&context{templatePath: path}, f)
}
