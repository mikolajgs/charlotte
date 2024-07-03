package shellstep

import (
	structvalidator "github.com/nicholasgasior/struct-validator"
)

type ShellStep struct {
	Name            string `yaml:"name" validation:"req lenmin:2 lenmax:80"`
	Description     string `yaml:"description" validation:"lenmax:240"`
	ContinueOnError bool   `yaml:"continue_on_error,omitempty"`
	Script          string `yaml:"script" validation:"req lenmin:1"`
	Outputs         map[string]string `yaml:"outputs"`
	RunScript       string
	ID              string `yaml:"id"`
}

func (s *ShellStep) GetName() string {
	return s.Name
}

func (s *ShellStep) GetDescription() string {
	return s.Description
}

func (s *ShellStep) GetContinueOnError() bool {
	return s.ContinueOnError
}

func (s *ShellStep) GetScript() string {
	return s.Script
}

func (s *ShellStep) GetRunScript() string {
	return s.RunScript
}

func (s *ShellStep) SetRunScript(sc string) {
	s.RunScript = sc
}

func (s *ShellStep) GetOutputs() map[string]string {
	return s.Outputs
}

func (s *ShellStep) SetOutput(n string, v string) {
	s.Outputs[n] = v
}

func (s *ShellStep) GetID() string {
	return s.ID
}

func (s *ShellStep) Validate() (bool, map[string]int) {
	isValid, failedFields := structvalidator.Validate(s, &structvalidator.ValidationOptions{})
	return isValid, failedFields
}
