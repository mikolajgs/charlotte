package input

import (
	structvalidator "github.com/mikolajgs/struct-validator"
	_ "gopkg.in/yaml.v2"
)

type Input struct {
	Name         string `yaml:"name" validation_regexp:"^[0-9a-zA-Z_]{2,50}$"`
	Description  string `yaml:"description" validation:"lenmax:240"`
	Required     bool   `yaml:"required"`
	Default      string `yaml:"default"`
	RegExp       string `yaml:"regexp"`
	RunValue     string `yaml:"value"`
}

func (i *Input) Validate() (bool, map[string]int) {
	isValid, failedFields := structvalidator.Validate(i, &structvalidator.ValidationOptions{})
	return isValid, failedFields
}
