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
	return runCommand(ctx, args[0], args[1:]...)
}

func shellAction(ctx *context, command interface{}) (interface{}, error) {
	return runCommand(ctx, "sh", "-c", command.(string))
}

func runCommand(ctx *context, command string, args ...string) (interface{}, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = filepath.Dir(ctx.templatePath)

	result, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.TrimSpace(string(result)), nil
}
