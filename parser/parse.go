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

type parser struct {
	context *context
}

func NewParser() *parser {
	return &parser{}
}

func (p *parser) ParseFile(path string) (interface{}, error) {
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

	p.context = &context{templatePath: path}
	defer func (contextPtr **context) { contextPtr = nil }(&p.context)

	return p.parse(f)
}

func (p *parser) parse(reader io.Reader) (interface{}, error) {
	var result interface{}
	if err := json.NewDecoder(reader).Decode(&result); err != nil {
		return nil, err
	}

	if err := p.walk(result); err != nil {
		return nil, err
	}

	return result, nil
}

// TODO: Verify there won't be problems with too deep recursion.
func (p *parser) walk(raw interface{}) error {
	switch node := raw.(type) {
	case []interface{}:
		for _, v := range node {
			if err := p.walk(v); err != nil {
				return err
			}
		}
	case map[string]interface{}:
		for k, v := range node {
			if transformation, ok := transformations[k]; ok {
				delete(node, k)

				data, err := transformation(p.context, v)
				if err != nil {
					return err
				}
				if data != nil {
					merge(node, data.(map[string]interface{}))
				}

				// Transformation changes the template, walk must be restarted.
				return p.walk(raw)
			}

			if err := p.walk(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func merge(dest map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		if _, ok := dest[k]; !ok {
			// Do not overwrite value if it is already present (mitigates the risk of accidentally overwriting
			// values defined in the containings template by values from included templates).
			dest[k] = v
		}
	}
}
