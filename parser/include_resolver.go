package parser

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"path/filepath"
)

type includeResolverConfig struct {
	// Base path for includes specified with relative path.
	basePath string
}

type includeResolver struct {
	config *includeResolverConfig

	result map[string]interface{}
}

func (resolver *includeResolver) resolve(include interface{}) error {
	switch path := include.(type) {
	case string:
		return resolver.includeFromFile(path)
	case []interface{}:
		return resolver.includeFromFiles(path)
	default:
		return errors.New(fmt.Sprintf("Trying to include '%v' which is neither a string nor an array", path))
	}
}

func (resolver *includeResolver) includeFromFiles(paths []interface{}) error {
	var errs error
	for _, path := range paths {
		pathStr, ok := path.(string)
		if !ok {
			errs = multierror.Append(errs, fmt.Errorf("Trying to include '%v' which is not a string", path))
			continue
		}
		if err := resolver.includeFromFile(pathStr); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

func (resolver *includeResolver) includeFromFile(path string) error {
	result, err := NewParser().ParseFile(resolver.resolvePath(path))
	if err != nil {
		return err
	}
	resolver.include(result.(map[string]interface{}))
	return nil
}

func (resolver *includeResolver) resolvePath(path string) string {
	if resolver.config.basePath != "" && !filepath.IsAbs(path) {
		return filepath.Join(resolver.config.basePath, path)
	}
	return path
}

func (resolver *includeResolver) include(src map[string]interface{}) {
	for srcKey, srcValue := range src {
		if _, ok := resolver.result[srcKey]; !ok {
			// Do not overwrite a value if it is already there. It makes ordering of includes important.
			resolver.result[srcKey] = srcValue
		}
	}
}
