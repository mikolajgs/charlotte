package job

import (
	"charlotte/pkg/input"
	"charlotte/pkg/output"
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	shellscriptstep "charlotte/pkg/step/shell-script"
	"fmt"
	"os"

	"github.com/mikogs/go-valifieldator"
	_ "gopkg.in/yaml.v2"
	"gopkg.in/yaml.v3"
)

type Job struct {
	Name        string           `yaml:"name" validation:"req lenmin:2 lenmax:80"`
	Description string           `yaml:"description" validation:"lenmax:240"`
	Inputs      []*input.Input   `yaml:"inputs"`
	Outputs     []*output.Output `yaml:"outputs"`
	Steps       interface{}      `yaml:"steps"`
}

func (j *Job) Validate() (bool, map[string]int) {
	isValid, failedFields := valifieldator.Validate(j, &valifieldator.ValidationOptions{})
	return isValid, failedFields
}

// Run executes a Step in a specific RuntimeEnvironment.
func (j *Job) Run(runtime runtimeenvironment.IRuntimeEnvironment) error {

	err := runtime.Create(j.Steps.([]step.IStep))
	if err != nil {
		return fmt.Errorf("error creating docker: %w", err)
	}

	defer runtime.Destroy(j.Steps.([]step.IStep))

	for i, step := range j.Steps.([]step.IStep) {
		fOut, fErr, err := runtime.Run(step, i)
		fmt.Fprintf(os.Stdout, "*** Step stdout file: %s\n", fOut)
		fmt.Fprintf(os.Stdout, "*** Step stderr file: %s\n", fErr)
		if err != nil && !step.GetContinueOnError() {
			return fmt.Errorf("step %s failed with: %w", step.GetName(), err)
		}
	}

	return nil
}

// NewFromFile creates Job instance from path to a YAML file.
func NewFromFile(f string) (*Job, error) {
	var j Job
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", f, err)
	}
	err = yaml.Unmarshal(b, &j)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal %s: %w", f, err)
	}

	newSteps := make([]step.IStep, 0)
	for i, s := range j.Steps.([]interface{}) {
		step := s.(map[string]interface{})
		typ, ok := step["type"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid type of step %s: %w", f, err)
		}
		if typ == "shellScript" {
			var newStep shellscriptstep.ShellScriptStep
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
	j.Steps = newSteps

	return &j, nil
}
