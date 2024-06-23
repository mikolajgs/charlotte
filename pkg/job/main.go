package job

import (
	"charlotte/pkg/input"
	jobrunresult "charlotte/pkg/job-run-result"
	"charlotte/pkg/output"
	runtimeenvironment "charlotte/pkg/runtime-environment"
	"charlotte/pkg/step"
	steprunresult "charlotte/pkg/step-run-result"
	shellscriptstep "charlotte/pkg/step/shell-script"
	"fmt"
	"os"

	structvalidator "github.com/nicholasgasior/struct-validator"
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
	isValid, failedFields := structvalidator.Validate(j, &structvalidator.ValidationOptions{})
	return isValid, failedFields
}

// Run executes a Step in a specific RuntimeEnvironment.
func (j *Job) Run(runtime runtimeenvironment.IRuntimeEnvironment) (*jobrunresult.JobRunResult) {
	jobRunResult := &jobrunresult.JobRunResult{}

	err := runtime.Create(j.Steps.([]step.IStep))
	if err != nil {
		jobRunResult.Error = fmt.Errorf("error creating runtimeenv: %w", err)
		jobRunResult.Success = false
		return jobRunResult
	}

	defer runtime.Destroy(j.Steps.([]step.IStep))

	jobRunResult.StepRunResults = make([]*steprunresult.StepRunResult, 0)
	jobRunResult.StepsWithErrors = make([]int, 0)

	for i, step := range j.Steps.([]step.IStep) {
		fOut, fErr, err := runtime.Run(step, i)
		//fmt.Fprintf(os.Stdout, "*** Step stdout file: %s\n", fOut)
		//fmt.Fprintf(os.Stdout, "*** Step stderr file: %s\n", fErr)

		suc := true
		if err != nil {
			suc = false
			jobRunResult.StepsWithErrors = append(jobRunResult.StepsWithErrors, i)
		}

		jobRunResult.StepRunResults = append(jobRunResult.StepRunResults, &steprunresult.StepRunResult{
			Success: suc,
			StderrFile: fErr,
			StdoutFile: fOut,
			Error: err,
		})

		if err != nil && !step.GetContinueOnError() {
			jobRunResult.Error = fmt.Errorf("step %s failed with: %w", step.GetName(), err)
			jobRunResult.Success = false
			jobRunResult.BreakingStep = i

			return jobRunResult
		}
	}

	return jobRunResult
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
