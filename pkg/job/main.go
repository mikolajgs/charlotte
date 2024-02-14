package job

import (
	"charlotte/pkg/input"
	"charlotte/pkg/output"
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	shellscriptstep "charlotte/pkg/step/shell-script"
	"fmt"
	"io/ioutil"
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
func (j *Job) Run(runtime runtimeenvironment.IRuntimeEnvironment) (int, []string) {
	for _, step := range j.Steps.([]step.IStep) {
		exitCode, errors, fOut, fErr := runtime.Run(step)
		fmt.Fprintf(os.Stdout, "exitCode: %d\n", exitCode)
		fmt.Fprintf(os.Stdout, "errors: %v\n", errors)
		fmt.Fprintf(os.Stdout, "fOut: %s\n", fOut)
		fmt.Fprintf(os.Stdout, "fErr: %s\n", fErr)
		if exitCode != 0 && !step.GetContinueOnError() {
			return exitCode, errors
		}
	}

	return 0, []string{}
}

// NewFromFile creates Job instance from path to a YAML file.
func NewFromFile(f string) (*Job, error) {
	var j Job
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("Cannot read file %s: %w\n", f, err)
	}
	err = yaml.Unmarshal(b, &j)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal %s: %w\n", f, err)
	}

	newSteps := make([]step.IStep, 0)
	for i, s := range j.Steps.([]interface{}) {
		step := s.(map[string]interface{})
		typ, ok := step["type"].(string)
		if !ok {
			return nil, fmt.Errorf("Invalid type of step %d: %w\n", f, err)
		}
		if typ == "shellScript" {
			var newStep shellscriptstep.ShellScriptStep
			stepBytes, err := yaml.Marshal(s)
			if err != nil {
				return nil, fmt.Errorf("Error marshalling step %d: %w\n", i, err)
			}
			err = yaml.Unmarshal(stepBytes, &newStep)
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling marshalled step %d: %w\n", i, err)
			}
			newSteps = append(newSteps, &newStep)
		}
	}
	j.Steps = newSteps

	return &j, nil
}
