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
	var err error
	switch node := raw.(type) {
	case []interface{}:
		err = p.walkArray(node)
	case map[string]interface{}:
		err = p.walkObject(node)
	}
	return err
}

func (p *parser) walkArray(array []interface{}) error {
	for i, v := range array {
		newValue, err := p.tryCallAction(v)
		if err != nil {
			return err
		} else if newValue != nil {
			array[i] = newValue
		}

		if err := p.walk(v); err != nil {
			return err
		}
	}
	return nil
}

func (p* parser) walkObject(object map[string]interface{}) error {
	for k, v := range object {
		if transformation, ok := transformations[k]; ok {
			delete(object, k)

			data, err := transformation(p.context, v)
			if err != nil {
				return err
			}
			if data != nil {
				merge(object, data.(map[string]interface{}))
			}

			// Transformation changes the template, walk must be restarted.
			return p.walkObject(object)
		}

		newValue, err := p.tryCallAction(v)
		if err != nil {
			return err
		} else if newValue != nil {
			object[k] = newValue
		}

		if err := p.walk(v); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) tryCallAction(value interface{}) (interface{}, error) {
	action := p.tryParseAction(value)
	if action == nil {
		return nil, nil
	}
	return p.callAction(action)
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
