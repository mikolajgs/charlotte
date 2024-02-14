package input

import (
	"github.com/mikogs/go-valifieldator"
	_ "gopkg.in/yaml.v2"
)

type Input struct {
	Name         string `yaml:"name" validation_regexp:"^[0-9a-zA-Z_]{2,50}$"`
	Description  string `yaml:"description" validation:"lenmax:240"`
	Required     bool   `yaml:"required"`
	DefaultValue string `yaml:"defaultValue"`
}

func (i *Input) Validate() (bool, map[string]int) {
	isValid, failedFields := valifieldator.Validate(i, &valifieldator.ValidationOptions{})
	return isValid, failedFields
}
