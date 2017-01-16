package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	defer func(contextPtr **context) { contextPtr = nil }(&p.context)

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
		for i, v := range node {
			action := p.tryParseAction(v)
			if action != nil {
				newValue, err := p.callAction(action)
				if err != nil {
					return err
				}
				node[i] = newValue
			}

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

			action := p.tryParseAction(v)
			if action != nil {
				newValue, err := p.callAction(action)
				if err != nil {
					return err
				}
				node[k] = newValue
			}

			if err := p.walk(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) tryParseAction(value interface{}) []string {
	switch value := value.(type) {
	case string:
		if strings.HasPrefix(value, "ppp-") {
			return strings.SplitN(value, ":", 2)
		}
	}
	return nil
}

func (p *parser) callAction(action []string) (interface{}, error) {
	name := action[0]
	body := getIndexOr(action, 1, "").(string)

	if handler, ok := actions[name]; ok {
		return handler(p.context, body)
	}

	return nil, fmt.Errorf("Action %s is not supported", name)
}
