package shellscriptstep

import (
	structvalidator "github.com/nicholasgasior/struct-validator"
)

type ShellScriptStep struct {
	Name            string `yaml:"name" validation:"req lenmin:2 lenmax:80"`
	Description     string `yaml:"description" validation:"lenmax:240"`
	ContinueOnError bool   `yaml:"continue_on_error,omitempty"`
	Script          string `yaml:"script" validation:"req lenmin:1"`
}

func (s *ShellScriptStep) GetName() string {
	return s.Name
}

func (s *ShellScriptStep) GetDescription() string {
	return s.Description
}

func (s *ShellScriptStep) GetContinueOnError() bool {
	return s.ContinueOnError
}

func (s *ShellScriptStep) GetScript() string {
	return s.Script
}

func (s *ShellScriptStep) Validate() (bool, map[string]int) {
	isValid, failedFields := structvalidator.Validate(s, &structvalidator.ValidationOptions{})
	return isValid, failedFields
}
