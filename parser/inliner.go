package parser

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"path/filepath"
)

type inlinerConfig struct {
	// Base path for includes specified with relative path.
	basePath string
}

type inliner struct {
	config *inlinerConfig

	result map[string]interface{}
}

func (i *inliner) inline(path interface{}) error {
	switch path := path.(type) {
	case string:
		return i.inlineFromFile(path)
	case []interface{}:
		return i.inlineFromFiles(path)
	default:
		return errors.New(fmt.Sprintf("Trying to include '%v' which is neither a string nor an array", path))
	}
}

func (i *inliner) inlineFromFiles(paths []interface{}) error {
	var errs error
	for _, path := range paths {
		pathStr, ok := path.(string)
		if !ok {
			errs = multierror.Append(errs, fmt.Errorf("Trying to include '%v' which is not a string", path))
			continue
		}
		if err := i.inlineFromFile(pathStr); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

func (i *inliner) inlineFromFile(path string) error {
	result, err := NewParser().ParseFile(i.resolvePath(path))
	if err != nil {
		return err
	}
	merge(i.result, result.(map[string]interface{}))
	return nil
}

func (i *inliner) resolvePath(path string) string {
	if i.config.basePath != "" && !filepath.IsAbs(path) {
		return filepath.Join(i.config.basePath, path)
	}
	return path
}
