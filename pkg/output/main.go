package output

import (
	structvalidator "github.com/nicholasgasior/struct-validator"
	_ "gopkg.in/yaml.v2"
)

type Output struct {
	Name         string `yaml:"name" validation_regexp:"^[0-9a-zA-Z_]{2,50}$"`
	Desc         string `yaml:"desc" validation:"lenmax:240"`
	DefaultValue string `yaml:"defaultValue"`
}

func (o *Output) Validate() (bool, map[string]int) {
	isValid, failedFields := structvalidator.Validate(o, &structvalidator.ValidationOptions{})
	return isValid, failedFields
}
