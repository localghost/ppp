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

func inlineTransformation(ctx *context, path interface{}) (interface{}, error) {
	inliner := &inliner{
		config: &inlinerConfig{basePath: filepath.Dir(ctx.templatePath)},
		result: make(map[string]interface{}),
	}

	if err := inliner.inline(path); err != nil {
		return nil, err
	}

	return inliner.result, nil
}
