package job

import (
	"charlotte/pkg/input"
	"charlotte/pkg/output"
	"charlotte/pkg/step"
	shellstep "charlotte/pkg/step/shell"
	"fmt"
	"os"

	_ "gopkg.in/yaml.v2"
	"gopkg.in/yaml.v3"
)

type Job struct {
	Name        string           `yaml:"name" validation:"req lenmin:2 lenmax:80"`
	Description string           `yaml:"description" validation:"lenmax:240"`
	Inputs      map[string]*input.Input   `yaml:"inputs"`
	Outputs     map[string]*output.Output `yaml:"outputs"`
	Steps       interface{}      `yaml:"steps"`
	Variables		map[string]string `yaml:"variables"`
	Environment map[string]string `yaml:"environment"`
}

// NewFromBytes creates Job instance from bytes array
func NewFromBytes(b []byte) (*Job, error) {
	var j Job

	err := yaml.Unmarshal(b, &j)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal: %w", err)
	}

	newSteps := make([]step.IStep, 0)
	for i, s := range j.Steps.([]interface{}) {
		step := s.(map[string]interface{})
		typ, ok := step["type"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid type of step %d: %w", i, err)
		}
		if typ == "shell" {
			var newStep shellstep.ShellStep
			stepBytes, err := yaml.Marshal(s)
			if err != nil {
				return nil, fmt.Errorf("error marshalling step %d: %w", i, err)
			}

			err = yaml.Unmarshal(stepBytes, &newStep)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling marshalled step %d: %w", i, err)
			}

			newSteps = append(newSteps, &newStep)
		}
	}

	for n := range j.Inputs {
		j.Inputs[n].Name = n
	}

	j.Steps = newSteps

	return &j, nil
}

// NewFromFile creates Job instance from path to a YAML file.
func NewFromFile(f string) (*Job, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", f, err)
	}
	j, err := NewFromBytes(b)
	if err != nil {
		return nil, fmt.Errorf("cannot create from bytes: %w", err)
	}
	return j, nil
}
