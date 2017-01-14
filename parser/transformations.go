package parser

import (
	"path/filepath"
)

type transformation func(*context, interface{}) (interface{}, error)

var transformations map[string]transformation

func init() {
	transformations = map[string]transformation{
		"ppp-inline": inlineTransformation,
	}
}

func inlineTransformation(ctx *context, data interface{}) (interface{}, error) {
	resolver := &includeResolver{
		config: &includeResolverConfig{basePath: filepath.Dir(ctx.templatePath)},
		result: make(map[string]interface{}),
	}

	if err := resolver.resolve(data); err != nil {
		return nil, err
	}

	return resolver.result, nil
}
