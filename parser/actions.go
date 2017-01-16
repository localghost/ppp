package parser

import (
	"path/filepath"
	"os/exec"
	"strings"
)

type action func(*context, interface{}) (interface{}, error)

var actions map[string]action

func init() {
	actions = map[string]action{
		"ppp-exec": execAction,
		"ppp-shell": shellAction,
	}
}

func execAction(ctx *context, command interface{}) (interface{}, error) {
	args := strings.Fields(command.(string))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = filepath.Dir(ctx.templatePath)

	result, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.TrimSpace(string(result)), nil
}

func shellAction(ctx *context, command interface{}) (interface{}, error) {
	cmd := exec.Command("sh", "-c", command.(string))
	cmd.Dir = filepath.Dir(ctx.templatePath)

	result, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.TrimSpace(string(result)), nil
}
